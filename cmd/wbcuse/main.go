package main

import (
	"os"

	"github.com/pojntfx/webpipe/pkg/cuse"
)

func main() {
	device := cuse.NewEchoDevice()

	if err := cuse.MountDevice(
		device,

		69,
		69,
		"wbcuse",

		os.Args,
	); err != nil {
		panic(err)
	}
}
