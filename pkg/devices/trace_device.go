package devices

import (
	"log"

	"github.com/pojntfx/webpipe/pkg/cuse"
)

type TraceDevice struct{}

func NewTraceDevice() *TraceDevice {
	return &TraceDevice{}
}

func (d *TraceDevice) Init(userdata cuse.Void, conn cuse.Conn) {
	log.Println("Init", userdata, conn)
}

func (d *TraceDevice) InitDone(userdata cuse.Void) {
	log.Println("InitDone", userdata)
}

func (d *TraceDevice) Destroy(userdata cuse.Void) {
	log.Println("Destroy", userdata)
}

func (d *TraceDevice) Open(req cuse.Request, fi cuse.FileInfo) {
	log.Println("Open", req, fi)

	if err := cuse.ReplyOpen(req, fi); err != nil {
		panic(err)
	}
}

func (d *TraceDevice) Read(req cuse.Request, size cuse.Size, off cuse.Offset, fi cuse.FileInfo) {
	log.Println("Read", req, size, off, fi)

	if err := cuse.ReplyBuf(req, []byte{}); err != nil {
		panic(err)
	}
}

func (d *TraceDevice) Write(req cuse.Request, buf cuse.Buffer, size cuse.Size, off cuse.Offset, fi cuse.FileInfo) {
	log.Println("Write", req, buf, size, off, fi)

	if err := cuse.ReplyWrite(req, int(cuse.SizeToUint64(size))); err != nil {
		panic(err)
	}
}

func (d *TraceDevice) Flush(req cuse.Request, fi cuse.FileInfo) {
	log.Println("Flush", req, fi)

	if err := cuse.ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *TraceDevice) Release(req cuse.Request, fi cuse.FileInfo) {
	log.Println("Release", req, fi)

	if err := cuse.ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *TraceDevice) Fsync(req cuse.Request, datasync int, fi cuse.FileInfo) {
	log.Println("Fsync", req, datasync, fi)

	if err := cuse.ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *TraceDevice) Ioctl(req cuse.Request, cmd int, arg cuse.Void, fi cuse.FileInfo, flags uint, inputBuf cuse.Void, inputBufSize cuse.Size, outBufSize cuse.Size) {
	log.Println("Ioctl", req, cmd, arg, fi, flags, inputBuf, inputBufSize, outBufSize)

	if err := cuse.ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *TraceDevice) Poll(req cuse.Request, fi cuse.FileInfo, ph cuse.PollHandle) {
	log.Println("Poll", req, fi, ph)

	if err := cuse.ReplyPoll(req, 0); err != nil {
		panic(err)
	}
}
