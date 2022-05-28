package main

import (
	"flag"
	"os"

	"github.com/pojntfx/webpipe/pkg/cuse"
	"github.com/pojntfx/webpipe/pkg/devices"
)

func main() {
	flag.Parse()

	device := devices.NewTraceDevice()

	if err := cuse.MountDevice(
		device,

		70,
		70,
		"wbcuse-trace",

		append([]string{os.Args[0]}, flag.Args()...),
	); err != nil {
		panic(err)
	}
}
