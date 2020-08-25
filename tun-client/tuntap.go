package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

func startTap() *water.Interface {
	fmt.Print(".")
	iface, err := water.NewTAP("scsi0")

	cmd := exec.Command("/usr/bin/ip", "link", "set", "up", "dev", "scsi0")
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()

	if err != nil {
		log.Fatalf("That's no fun. Can't make a tap device: %v", err)
	}

	return iface

	// go func() {
	// 	for pkt2 := range inboundPackets {
	// 		if len(pkt2) != 0 {
	// 			if pkt2[12] != 0x00 {
	// 				if *debugEnabled {
	// 					fmt.Print(">")
	// 				}
	// 				_, err := iface.Write(pkt2)
	// 				if err != nil {
	// 					log.Fatalf("Can't write to tap device, I don't know how this happens but its likely fatal: %v", err)
	// 				}
	// 			}
	// 		}
	// 	}
	// }()

	// for {
	// 	pkt := make([]byte, 512*3)
	// 	for i := 0; i < len(pkt); i++ {
	// 		pkt[i] = 0x00
	// 	}
	// 	n, err := iface.Read(pkt)
	// 	if err != nil {
	// 		log.Fatalf("Can't read from tap device, I don't know how this happens but its likely fatal: %v", err)
	// 	}

	// 	outboundPackets <- pkt[:n]
	// 	if *debugEnabled {
	// 		fmt.Print("<")
	// 	}
	// }
}

var (
// outboundPackets chan []byte
// inboundPackets  chan []byte
)
