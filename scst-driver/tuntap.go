package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

func startTap() *water.Interface {
	// outboundPackets = make(chan []byte, 3)
	// inboundPackets = make(chan []byte, 3)

	fmt.Print(".")
	iface, err := water.NewTAP("scsi0")
	if err != nil {
		log.Fatalf("That's no fun. Can't make a tap device: %v", err)
	}

	cmd := exec.Command("/usr/bin/ip", "link", "set", "up", "dev", "scsi0")
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()

	return iface
}

var (
	outboundPackets chan []byte
	inboundPackets  chan []byte
)
