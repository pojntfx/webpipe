package main

import (
	"flag"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	src := flag.String("src", "src.entangled", "Path in which to create the source pipe")
	dst := flag.String("dst", "dst.entangled", "Path in which to create the destination pipe")

	flag.Parse()

	_ = os.Remove(*src)
	_ = os.Remove(*dst)

	if err := unix.Mkfifo(*src, 0666); err != nil {
		panic(err)
	}

	if err := unix.Mkfifo(*dst, 0666); err != nil {
		panic(err)
	}

	srcFile, err := os.OpenFile(*src, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}

	dstFile, err := os.OpenFile(*dst, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}

	errs := make(chan error)

	go func() {
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			errs <- err
		}
	}()

	go func() {
		if _, err := io.Copy(srcFile, dstFile); err != nil {
			errs <- err
		}
	}()

	if err := <-errs; err != nil {
		panic(err)
	}
}
