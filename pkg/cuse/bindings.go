package cuse

// #cgo CFLAGS: -Wno-error=implicit-function-declaration
// #cgo pkg-config: fuse3
// #include "stdlib.h"
// #include "cuse.h"
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/mattn/go-pointer"
)

type Void *C.void
type Conn *C.fuse_conn_info
type Request *C.struct_fuse_req
type FileInfo *C.fuse_file_info
type Size C.ulong
type Offset C.off_t
type Buffer *C.char
type PollHandle *C.fuse_pollhandle

func ReplyOpen(req Request, fi FileInfo) error {
	if ret := C.fuse_reply_open(req, fi); ret != 0 {
		return fmt.Errorf("could not reply with open: error code %v", ret)
	}

	return nil
}

func ReplyBuf(req Request, buf []byte) error {
	p := C.CString(string(buf))
	defer C.free(unsafe.Pointer(p))

	if ret := C.fuse_reply_buf(req, p, C.ulong(len(buf))); ret != 0 {
		return fmt.Errorf("could not reply with buffer: error code %v", ret)
	}

	return nil
}

func ReplyWrite(req Request, n int) error {
	if ret := C.fuse_reply_write(req, C.ulong(n)); ret != 0 {
		return fmt.Errorf("could not reply with write: error code %v", ret)
	}

	return nil
}

func ReplyError(req Request, err int) error {
	if ret := C.fuse_reply_err(req, C.int(err)); ret != 0 {
		return fmt.Errorf("could not reply with error: error code %v", ret)
	}

	return nil
}

func ReplyPoll(req Request, revents int) error {
	if ret := C.fuse_reply_poll(req, C.uint(revents)); ret != 0 {
		return fmt.Errorf("could not reply with poll: error code %v", ret)
	}

	return nil
}

func BufferToBytes(buf Buffer) []byte {
	return C.GoBytes(unsafe.Pointer(buf), C.int(unsafe.Sizeof(buf)))
}

func OffsetToInt64(off Offset) int64 {
	return int64(C.long(off))
}

func SizeToUint64(size Size) uint64 {
	return uint64(size)
}

//export wbcuse_init
func wbcuse_init(device unsafe.Pointer, userdata Void, conn Conn) {
	(pointer.Restore(device)).(*deviceContainer).device.Init(userdata, conn)
}

//export wbcuse_init_done
func wbcuse_init_done(device unsafe.Pointer, userdata Void) {
	(pointer.Restore(device)).(*deviceContainer).device.InitDone(userdata)
}

//export wbcuse_destroy
func wbcuse_destroy(device unsafe.Pointer, userdata Void) {
	(pointer.Restore(device)).(*deviceContainer).device.Destroy(userdata)
}

//export wbcuse_open
func wbcuse_open(device unsafe.Pointer, req Request, fi FileInfo) {
	(pointer.Restore(device)).(*deviceContainer).device.Open(req, fi)
}

//export wbcuse_read
func wbcuse_read(device unsafe.Pointer, req Request, size Size, off Offset, fi FileInfo) {
	(pointer.Restore(device)).(*deviceContainer).device.Read(req, size, off, fi)
}

//export wbcuse_write
func wbcuse_write(device unsafe.Pointer, req Request, buf Buffer, size Size, off Offset, fi FileInfo) {
	(pointer.Restore(device)).(*deviceContainer).device.Write(req, buf, size, off, fi)
}

//export wbcuse_flush
func wbcuse_flush(device unsafe.Pointer, req Request, fi FileInfo) {
	(pointer.Restore(device)).(*deviceContainer).device.Flush(req, fi)
}

//export wbcuse_release
func wbcuse_release(device unsafe.Pointer, req Request, fi FileInfo) {
	(pointer.Restore(device)).(*deviceContainer).device.Release(req, fi)
}

//export wbcuse_fsync
func wbcuse_fsync(device unsafe.Pointer, req Request, datasync C.int, fi FileInfo) {
	(pointer.Restore(device)).(*deviceContainer).device.Fsync(req, int(datasync), fi)
}

//export wbcuse_ioctl
func wbcuse_ioctl(device unsafe.Pointer, req Request, cmd C.int, arg Void, fi FileInfo, flags C.uint, in_buf Void, in_bufz Size, out_bufsz Size) {
	(pointer.Restore(device)).(*deviceContainer).device.Ioctl(req, int(cmd), arg, fi, uint(flags), in_buf, in_bufz, out_bufsz)
}

//export wbcuse_poll
func wbcuse_poll(device unsafe.Pointer, req Request, fi FileInfo, ph PollHandle) {
	(pointer.Restore(device)).(*deviceContainer).device.Poll(req, fi, ph)
}

func MountDevice(
	device Device,

	major uint,
	minor uint,
	name string,

	fuseArgs []string,
) error {
	cargs := []*C.char{}
	for _, arg := range fuseArgs {
		cargs = append(cargs, C.CString(arg))
	}

	if ret := C.wbcuse_start(
		pointer.Save(&deviceContainer{device}),

		C.uint(major),
		C.uint(minor),
		C.CString("DEVNAME="+name),

		C.int(len(cargs)),
		&cargs[0],
	); ret != 0 {
		return fmt.Errorf("could not start CUSE device: %v", ret)
	}

	return nil
}
