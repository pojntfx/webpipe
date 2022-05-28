package main

import (
	"flag"
	"os"

	"github.com/pojntfx/webpipe/pkg/cuse"
	"github.com/pojntfx/webpipe/pkg/devices"
)

func main() {
	backendFlag := flag.String("backend", "/tmp/wbcuse.entangled", "Name of the file to use as the backend")

	flag.Parse()

	backend, err := os.OpenFile(*backendFlag, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	device := devices.NewFileDevice(backend)

	if err := cuse.MountDevice(
		device,

		69,
		69,
		"wbcuse-file",

		append([]string{os.Args[0]}, flag.Args()...),
	); err != nil {
		panic(err)
	}
}
