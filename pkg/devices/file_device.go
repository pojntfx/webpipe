package devices

import (
	"io"
	"log"

	"github.com/pojntfx/webpipe/pkg/cuse"
)

type readSeekerWriter interface {
	io.ReadSeeker
	io.WriterAt
}

type FileDevice struct {
	backend readSeekerWriter
}

func NewFileDevice(backend readSeekerWriter) *FileDevice {
	return &FileDevice{backend}
}

func (d *FileDevice) Init(userdata cuse.Void, conn cuse.Conn) {
	log.Println("Init", userdata, conn)
}

func (d *FileDevice) InitDone(userdata cuse.Void) {
	log.Println("InitDone", userdata)
}

func (d *FileDevice) Destroy(userdata cuse.Void) {
	log.Println("Destroy", userdata)
}

func (d *FileDevice) Open(req cuse.Request, fi cuse.FileInfo) {
	log.Println("Open", req, fi)

	if err := cuse.ReplyOpen(req, fi); err != nil {
		panic(err)
	}
}

func (d *FileDevice) Read(req cuse.Request, size cuse.Size, off cuse.Offset, fi cuse.FileInfo) {
	log.Println("Read", req, size, off, fi)

	if _, err := d.backend.Seek(cuse.OffsetToInt64(off), io.SeekStart); err != nil {
		panic(err)
	}

	buf := make([]byte, cuse.SizeToUint64(size))

	n, err := io.ReadFull(d.backend, buf)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			panic(err)
		}
	}

	if err := cuse.ReplyBuf(req, buf[:n]); err != nil {
		panic(err)
	}
}

func (d *FileDevice) Write(req cuse.Request, buf cuse.Buffer, size cuse.Size, off cuse.Offset, fi cuse.FileInfo) {
	log.Println("Write", req, buf, size, off, fi)

	_, err := d.backend.WriteAt(cuse.BufferToBytes(buf), cuse.OffsetToInt64(off))
	if err != nil {
		panic(err)
	}

	if err := cuse.ReplyWrite(req, int(cuse.SizeToUint64(size))); err != nil {
		panic(err)
	}
}

func (d *FileDevice) Flush(req cuse.Request, fi cuse.FileInfo) {
	log.Println("Flush", req, fi)

	if err := cuse.ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *FileDevice) Release(req cuse.Request, fi cuse.FileInfo) {
	log.Println("Release", req, fi)

	if err := cuse.ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *FileDevice) Fsync(req cuse.Request, datasync int, fi cuse.FileInfo) {
	log.Println("Fsync", req, datasync, fi)

	if err := cuse.ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *FileDevice) Ioctl(req cuse.Request, cmd int, arg cuse.Void, fi cuse.FileInfo, flags uint, inputBuf cuse.Void, inputBufSize cuse.Size, outBufSize cuse.Size) {
	log.Println("Ioctl", req, cmd, arg, fi, flags, inputBuf, inputBufSize, outBufSize)

	if err := cuse.ReplyError(req, 0); err != nil {
		panic(err)
	}
}

func (d *FileDevice) Poll(req cuse.Request, fi cuse.FileInfo, ph cuse.PollHandle) {
	log.Println("Poll", req, fi, ph)

	if err := cuse.ReplyPoll(req, 0); err != nil {
		panic(err)
	}
}
