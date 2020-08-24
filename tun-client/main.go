package main

import (
	"crypto/rand"
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
	f1, err := os.Open("/dev/sg1")
	if err != nil {
		log.Fatalf("Failed to open ATA device")
	}

	f2, err := os.Open("/dev/sg2")
	if err != nil {
		log.Fatalf("Failed to open ATA device")
	}

	go startTap()
	// sendReadSgio(f1)
	// sendReadSgio(f1)

	go func() {
		for {
			pkt2 := <-outboundPackets
			err := sendSgio(f1, pkt2)
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
				inboundPackets <- pkt
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
	var testbuf [1536]uint8
	inqCmdBlk[0] = 0x8a
	// inqCmdBlk[9] = 0xFF

	randLBA := make([]byte, 3)
	rand.Read(randLBA)
	copy(inqCmdBlk[6:], randLBA)

	// inqCmdBlk[12] = 0x05
	inqCmdBlk[13] = 0x03 // 512 (block size) * 3

	copy(testbuf[:], pkt)
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
		DxferLen:       1536,
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
	var testbuf [1536]uint8
	pkt = make([]byte, 1536)
	inqCmdBlk[0] = 0x88

	randLBA := make([]byte, 3)
	rand.Read(randLBA)
	copy(inqCmdBlk[6:], randLBA)

	inqCmdBlk[9] = 0xFF

	// inqCmdBlk[12] = 0x05
	inqCmdBlk[13] = 0x03 // 512 (block size) * 3

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
		DxferLen:       1536,
		Dxferp:         &testbuf[0],
	}

	if err := sgio.SgioSyscall(f, ioHdr); err != nil {
		return pkt, err
	}

	if err := sgio.CheckSense(ioHdr, &senseBuf); err != nil {
		return pkt, err
	}

	copy(pkt, testbuf[:])

	if *debugEnabled {
		if pkt[12] != 0 {
			// log.Printf("READ READ REEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE %#v", pkt)
			fmt.Print(",")
		} else {
			fmt.Print(".")
		}
	}
	return pkt[:], nil
}
