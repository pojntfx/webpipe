package cuse

// #cgo CFLAGS: -Wno-error=implicit-function-declaration
// #cgo pkg-config: fuse3
// #include "cuse.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type Void *C.void
type Conn *C.fuse_conn_info
type Request *C.fuse_req_t
type FileInfo *C.fuse_file_info
type Size *C.size_t
type Offset C.off_t
type Buffer *C.char
type PollHandle *C.fuse_pollhandle

//export wbcuse_init
func wbcuse_init(device unsafe.Pointer, userdata Void, conn Conn) {
	(*deviceContainer)(device).device.Init(userdata, conn)
}

//export wbcuse_init_done
func wbcuse_init_done(device unsafe.Pointer, userdata Void) {
	(*deviceContainer)(device).device.InitDone(userdata)
}

//export wbcuse_destroy
func wbcuse_destroy(device unsafe.Pointer, userdata Void) {
	(*deviceContainer)(device).device.Destroy(userdata)
}

//export wbcuse_open
func wbcuse_open(device unsafe.Pointer, req Request, fi FileInfo) {
	(*deviceContainer)(device).device.Open(req, fi)
}

//export wbcuse_read
func wbcuse_read(device unsafe.Pointer, req Request, size Size, off Offset, fi FileInfo) {
	(*deviceContainer)(device).device.Read(req, size, off, fi)
}

//export wbcuse_write
func wbcuse_write(device unsafe.Pointer, req Request, buf Buffer, size Size, off Offset, fi FileInfo) {
	(*deviceContainer)(device).device.Write(req, buf, size, off, fi)
}

//export wbcuse_flush
func wbcuse_flush(device unsafe.Pointer, req Request, fi FileInfo) {
	(*deviceContainer)(device).device.Flush(req, fi)
}

//export wbcuse_release
func wbcuse_release(device unsafe.Pointer, req Request, fi FileInfo) {
	(*deviceContainer)(device).device.Release(req, fi)
}

//export wbcuse_fsync
func wbcuse_fsync(device unsafe.Pointer, req Request, datasync C.int, fi FileInfo) {
	(*deviceContainer)(device).device.Fsync(req, int(datasync), fi)
}

//export wbcuse_ioctl
func wbcuse_ioctl(device unsafe.Pointer, req Request, cmd C.int, arg Void, fi FileInfo, flags C.uint, in_buf Void, in_bufz Size, out_bufsz Size) {
	(*deviceContainer)(device).device.Ioctl(req, int(cmd), arg, fi, uint(flags), in_buf, in_bufz, out_bufsz)
}

//export wbcuse_poll
func wbcuse_poll(device unsafe.Pointer, req Request, fi FileInfo, ph PollHandle) {
	(*deviceContainer)(device).device.Poll(req, fi, ph)
}

func OpenDevice(device Device, args []string) error {
	cargs := []*C.char{}
	for _, arg := range args {
		cargs = append(cargs, C.CString(arg))
	}

	if ret := C.wbcuse_start(unsafe.Pointer(&deviceContainer{device}), C.int(len(cargs)), &cargs[0]); ret != 0 {
		return fmt.Errorf("could not start CUSE device: %v", ret)
	}

	return nil
}
