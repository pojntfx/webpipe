package cuse

import (
	"io"
	"log"
)

type readSeekerWriter interface {
	io.ReadSeeker
	io.WriterAt
}

type EchoDevice struct {
	backend readSeekerWriter
}

func NewEchoDevice(backend readSeekerWriter) *EchoDevice {
	return &EchoDevice{backend}
}

func (d *EchoDevice) Init(userdata Void, conn Conn) {
	log.Println("Init", userdata, conn)
}

func (d *EchoDevice) InitDone(userdata Void) {
	log.Println("InitDone", userdata)
}

func (d *EchoDevice) Destroy(userdata Void) {
	log.Println("Destroy", userdata)
}

func (d *EchoDevice) Open(req Request, fi FileInfo) {
	log.Println("Open", req, fi)

	if err := ReplyOpen(req, fi); err != nil {
		panic(err)
	}
}

func (d *EchoDevice) Read(req Request, size Size, off Offset, fi FileInfo) {
	log.Println("Read", req, size, off, fi)

	if _, err := d.backend.Seek(OffsetToInt64(off), io.SeekStart); err != nil {
		panic(err)
	}

	buf := make([]byte, SizeToUint64(size))

	n, err := io.ReadFull(d.backend, buf)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			panic(err)
		}
	}

	if err := ReplyBuf(req, buf[:n]); err != nil {
		panic(err)
	}
}

func (d *EchoDevice) Write(req Request, buf Buffer, size Size, off Offset, fi FileInfo) {
	log.Println("Write", req, buf, size, off, fi)

	_, err := d.backend.WriteAt(BufferToBytes(buf), OffsetToInt64(off))
	if err != nil {
		panic(err)
	}

	if err := ReplyWrite(req, int(SizeToUint64(size))); err != nil {
		panic(err)
	}
}

func (d *EchoDevice) Flush(req Request, fi FileInfo) {
	log.Println("Flush", req, fi)
}

func (d *EchoDevice) Release(req Request, fi FileInfo) {
	log.Println("Release", req, fi)

	if err := ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *EchoDevice) Fsync(req Request, datasync int, fi FileInfo) {
	log.Println("Fsync", req, datasync, fi)

	if err := ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *EchoDevice) Ioctl(req Request, cmd int, arg Void, fi FileInfo, flags uint, inputBuf Void, inputBufSize Size, outBufSize Size) {
	log.Println("Ioctl", req, cmd, arg, fi, flags, inputBuf, inputBufSize, outBufSize)

	if err := ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *EchoDevice) Poll(req Request, fi FileInfo, ph PollHandle) {
	log.Println("Poll", req, fi, ph)

	if err := ReplyError(req, 0); err != nil {
		panic(err)
	}
}
