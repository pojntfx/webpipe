package main

import (
	"os"

	"github.com/pojntfx/webpipe/pkg/cuse"
)

func main() {
	device := cuse.NewEchoDevice()

	if err := cuse.OpenDevice(device, os.Args); err != nil {
		panic(err)
	}
}
