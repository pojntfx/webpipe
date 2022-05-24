package main

import (
	"os"

	"github.com/pojntfx/webpipe/pkg/cuse"
)

func main() {
	id := 0
	device := cuse.NewEchoDevice(id)

	cuse.GlobalRegistry.AddDevice(id, device)
	defer cuse.GlobalRegistry.RemoveDevice(id)

	if err := device.Open(os.Args); err != nil {
		panic(err)
	}
}
