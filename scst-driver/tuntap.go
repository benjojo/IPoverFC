package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

func startTap() {
	outboundPackets = make(chan []byte, 1)
	inboundPackets = make(chan []byte, 1)

	fmt.Print(".")
	iface, err := water.NewTAP("scsi0")
	if err != nil {
		log.Fatalf("That's no fun. Can't make a tap device: %v", err)
	}

	cmd := exec.Command("/usr/bin/ip", "link", "set", "up", "dev", "scsi0")
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()

	go func() {
		for pkt := range inboundPackets {
			if pkt[12] != 0x00 {
				_, err := iface.Write(pkt[:])
				if err != nil {
					log.Fatalf("Can't write to tap device, I don't know how this happens but its likely fatal: %v", err)
				}
			}
		}
	}()

	for {
		pkt := make([]byte, 512*3)
		n, err := iface.Read(pkt[:])
		if err != nil {
			log.Fatalf("Can't read from tap device, I don't know how this happens but its likely fatal: %v", err)
		}

		outboundPackets <- pkt[:n]
	}
}

var (
	outboundPackets chan []byte
	inboundPackets  chan []byte
)
