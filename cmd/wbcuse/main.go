package main

import (
	"flag"
	"os"

	"github.com/pojntfx/webpipe/pkg/cuse"
)

func main() {
	backendFlag := flag.String("backend", "/tmp/wbcuse.entangled", "Name of the file to use as the backend")

	flag.Parse()

	backend, err := os.OpenFile(*backendFlag, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	device := cuse.NewEchoDevice(backend)

	if err := cuse.MountDevice(
		device,

		69,
		69,
		"wbcuse",

		append([]string{os.Args[0]}, flag.Args()...),
	); err != nil {
		panic(err)
	}
}
