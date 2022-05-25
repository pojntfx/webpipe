package cuse

// #cgo CFLAGS: -Wno-error=implicit-function-declaration
// #cgo pkg-config: fuse3
// #include "cuse.h"
import "C"
import "fmt"

type Void *C.void
type Conn *C.fuse_conn_info
type Request *C.fuse_req_t
type FileInfo *C.fuse_file_info
type Size *C.size_t
type Offset C.off_t
type Buffer *C.char
type PollHandle *C.fuse_pollhandle

func getDeviceOrPanic(registry_id C.int) Device {
	device, err := GlobalRegistry.GetDevice(int(registry_id))
	if err != nil {
		panic(err)
	}

	return device
}

//export wbcuse_init
func wbcuse_init(registry_id C.int, userdata Void, conn Conn) {
	getDeviceOrPanic(registry_id).Init(userdata, conn)
}

//export wbcuse_init_done
func wbcuse_init_done(registry_id C.int, userdata Void) {
	getDeviceOrPanic(registry_id).InitDone(userdata)
}

//export wbcuse_destroy
func wbcuse_destroy(registry_id C.int, userdata Void) {
	getDeviceOrPanic(registry_id).Destroy(userdata)
}

//export wbcuse_open
func wbcuse_open(registry_id C.int, req Request, fi FileInfo) {
	getDeviceOrPanic(registry_id).Open(req, fi)
}

//export wbcuse_read
func wbcuse_read(registry_id C.int, req Request, size Size, off Offset, fi FileInfo) {
	getDeviceOrPanic(registry_id).Read(req, size, off, fi)
}

//export wbcuse_write
func wbcuse_write(registry_id C.int, req Request, buf Buffer, size Size, off Offset, fi FileInfo) {
	getDeviceOrPanic(registry_id).Write(req, buf, size, off, fi)
}

//export wbcuse_flush
func wbcuse_flush(registry_id C.int, req Request, fi FileInfo) {
	getDeviceOrPanic(registry_id).Flush(req, fi)
}

//export wbcuse_release
func wbcuse_release(registry_id C.int, req Request, fi FileInfo) {
	getDeviceOrPanic(registry_id).Release(req, fi)
}

//export wbcuse_fsync
func wbcuse_fsync(registry_id C.int, req Request, datasync C.int, fi FileInfo) {
	getDeviceOrPanic(registry_id).Fsync(req, int(datasync), fi)
}

//export wbcuse_ioctl
func wbcuse_ioctl(registry_id C.int, req Request, cmd C.int, arg Void, fi FileInfo, flags C.uint, in_buf Void, in_bufz Size, out_bufsz Size) {
	getDeviceOrPanic(registry_id).Ioctl(req, int(cmd), arg, fi, uint(flags), in_buf, in_bufz, out_bufsz)
}

//export wbcuse_poll
func wbcuse_poll(registry_id C.int, req Request, fi FileInfo, ph PollHandle) {
	getDeviceOrPanic(registry_id).Poll(req, fi, ph)
}

func StartCUSE(registryID int, args []string) error {
	cargs := []*C.char{}
	for _, arg := range args {
		cargs = append(cargs, C.CString(arg))
	}

	if ret := C.wbcuse_start(C.int(registryID), C.int(len(cargs)), &cargs[0]); ret != 0 {
		return fmt.Errorf("could not start CUSE device: %v", ret)
	}

	return nil
}
