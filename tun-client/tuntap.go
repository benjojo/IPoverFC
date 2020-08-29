package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

func startTap() *water.Interface {
	iface, err := water.NewTAP("scsi0")

	cmd := exec.Command("/usr/bin/ip", "link", "set", "up", "dev", "scsi0")
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()

	if err != nil {
		log.Fatalf("That's no fun. Can't make a tap device: %v", err)
	}

	return iface
}
