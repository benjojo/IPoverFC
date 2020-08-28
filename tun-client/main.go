package main

import (
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/benmcclelland/sgio"
)

var lastRead = 0
var debugEnabled = flag.Bool("debug", false, "Enable debug text")

func main() {
	flag.Parse()
	f1, err := os.Open("/dev/sg1")
	if err != nil {
		log.Fatalf("Failed to open ATA device")
	}

	f2, err := os.Open("/dev/sg2")
	if err != nil {
		log.Fatalf("Failed to open ATA device")
	}

	tuntap := startTap()
	// sendReadSgio(f1)
	// sendReadSgio(f1)

	go func() {
		pkt2 := make([]byte, 9000)
		for {
			// pkt2 := <-outboundPackets

			n, err := tuntap.Read(pkt2)
			if err != nil {
				log.Printf("TUNTAP ERROR ON READ: %v", err)
				time.Sleep(time.Second)
				continue
			}
			err = sendSgio(f1, pkt2[:n])
			if err != nil {
				log.Printf("ATA error on write %v", err)
				time.Sleep(time.Second)
			}
		}
	}()

	for {
		// hadRead := false
		pkt, err := sendReadSgio(f2)
		if err != nil {
			log.Printf("ATA error on read %v", err)
			time.Sleep(time.Second)
		} else {
			if len(pkt) != 0 {
				if pkt[12] != 0x00 {
					tuntap.Write(pkt)
				}
				// inboundPackets <- pkt
			}
		}
	}
}

const (
	sgAta16    = 0x85
	sgAta16Len = 16

	sgDxferNone = -1

	sgAtaProtoNonData = 3 << 1
	sgCdb2CheckCond   = 1 << 5
	ataUsingLba       = 1 << 6

	ataOpStandbyNow1 = 0xe0 // https://wiki.osdev.org/ATA/ATAPI_Power_Management
	ataOpStandbyNow2 = 0x94 // Retired in ATA4. Did not coexist with ATAPI.
)

/*
	SG_DXFER_NONE (-1)      // e.g. a SCSI Test Unit Ready command
	SG_DXFER_TO_DEV (-2)    // e.g. a SCSI WRITE command
	SG_DXFER_FROM_DEV (-3)  // e.g. a SCSI READ command
	SG_DXFER_TO_FROM_DEV (-4) // treated like SG_DXFER_FROM_DEV with t
					additional property than during indirect
					IO the user buffer is copied into the
					kernel buffers before the transfer
	SG_DXFER_UNKNOWN (-5)   // Unknown data direction
*/

func sendSgio(f *os.File, pkt []byte) error {
	// log.Printf("Packet to be ATA written %v", pkt)
	var inqCmdBlk [sgAta16Len]uint8
	var testbuf [19 * 512]uint8
	inqCmdBlk[0] = 0x8a
	// inqCmdBlk[9] = 0xFF

	randLBA := make([]byte, 3)
	rand.Read(randLBA)
	copy(inqCmdBlk[6:], randLBA)

	// inqCmdBlk[12] = 0x05
	// inqCmdBlk[13] = 0x03 // 512 (block size) * 3
	inqCmdBlk[13] = 0x13 // 512 (block size) * 19

	var PLenbytes [2]byte
	binary.BigEndian.PutUint16(PLenbytes[:], uint16(len(pkt)))
	testbuf[0] = PLenbytes[0]
	testbuf[1] = PLenbytes[1]

	copy(testbuf[2:], pkt)
	for i := 0; i < len(pkt); i++ {
		pkt[i] = 0x00
	}

	senseBuf := make([]byte, sgio.SENSE_BUF_LEN)
	ioHdr := &sgio.SgIoHdr{
		InterfaceID:    'S',
		DxferDirection: -2,
		Cmdp:           &inqCmdBlk[0],
		CmdLen:         sgAta16Len,
		Sbp:            &senseBuf[0],
		MxSbLen:        sgio.SENSE_BUF_LEN,
		Timeout:        0,
		DxferLen:       19 * 512,
		Dxferp:         &testbuf[0],
	}

	// log.Printf("Test %#v", testbuf)

	if err := sgio.SgioSyscall(f, ioHdr); err != nil {
		return err
	}

	if err := sgio.CheckSense(ioHdr, &senseBuf); err != nil {
		return err
	}
	return nil
}

func sendReadSgio(f *os.File) (pkt []byte, err error) {

	var inqCmdBlk [sgAta16Len]uint8
	// var testbuf [512 * 3]uint8
	var testbuf [512 * 19]uint8
	pkt = make([]byte, 512*64)
	inqCmdBlk[0] = 0x88

	randLBA := make([]byte, 3)
	rand.Read(randLBA)
	copy(inqCmdBlk[6:], randLBA)

	inqCmdBlk[9] = 0xFF

	// inqCmdBlk[12] = 0x05
	// inqCmdBlk[13] = 0x03 // 512 (block size) * 3
	inqCmdBlk[13] = 0x13 // 512 (block size) * 19 = 9k

	testbuf[0] = 0x00

	senseBuf := make([]byte, sgio.SENSE_BUF_LEN)
	ioHdr := &sgio.SgIoHdr{
		InterfaceID:    'S',
		DxferDirection: -3,
		Cmdp:           &inqCmdBlk[0],
		CmdLen:         sgAta16Len,
		Sbp:            &senseBuf[0],
		MxSbLen:        sgio.SENSE_BUF_LEN,
		Timeout:        0,
		DxferLen:       512 * 19,
		Dxferp:         &testbuf[0],
	}

	if err := sgio.SgioSyscall(f, ioHdr); err != nil {
		return pkt, err
	}

	if err := sgio.CheckSense(ioHdr, &senseBuf); err != nil {
		return pkt, err
	}

	copy(pkt, testbuf[:])
	if len(pkt) != 0 {
		PktLen := binary.BigEndian.Uint16(pkt[:2])
		if PktLen != 0 {
			if *debugEnabled {
				if pkt[12] != 0 {
					// log.Printf("READ READ REEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE %#v", pkt)
					fmt.Print(",")
				} else {
					fmt.Print(".")
				}
			}
			return pkt[2 : 2+PktLen], nil

		}
	}
	return pkt, nil
}
