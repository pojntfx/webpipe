package main

// #cgo CFLAGS: -Wno-error=implicit-function-declaration
// #cgo pkg-config: fuse3
// #include "cuse.h"
import "C"

import (
	"fmt"
	"os"
)

//export wbcuse_init
func wbcuse_init(userdata *C.void, conn *C.fuse_conn_info) {
	fmt.Println("wbcuse_init", userdata, conn)
}

//export wbcuse_init_done
func wbcuse_init_done(userdata *C.void) {
	fmt.Println("wbcuse_init_done", userdata)
}

//export wbcuse_destroy
func wbcuse_destroy(userdata *C.void) {
	fmt.Println("wbcuse_destroy", userdata)
}

//export wbcuse_open
func wbcuse_open(req *C.struct_fuse_req, fi *C.fuse_file_info) {
	fmt.Println("wbcuse_open", req, fi)

	C.fuse_reply_open(req, fi)
}

//export wbcuse_read
func wbcuse_read(req *C.fuse_req_t, size C.size_t, off C.off_t, fi *C.fuse_file_info) {
	fmt.Println("wbcuse_read", req, size, off, fi)
}

//export wbcuse_write
func wbcuse_write(req *C.fuse_req_t, buf *C.char, size C.size_t, off C.off_t, fi *C.fuse_file_info) {
	fmt.Println("wbcuse_write", req, buf, size, off, fi)
}

//export wbcuse_flush
func wbcuse_flush(req *C.fuse_req_t, fi *C.fuse_file_info) {
	fmt.Println("wbcuse_flush", req, fi)
}

//export wbcuse_release
func wbcuse_release(req *C.fuse_req_t, fi *C.fuse_file_info) {
	fmt.Println("wbcuse_release", req, fi)
}

//export wbcuse_fsync
func wbcuse_fsync(req *C.fuse_req_t, datasync C.int, fi *C.fuse_file_info) {
	fmt.Println("wbcuse_fsync", req, datasync, fi)
}

//export wbcuse_ioctl
func wbcuse_ioctl(req *C.fuse_req_t, cmd C.int, arg *C.void, fi *C.fuse_file_info, flags C.uint, in_buf *C.void, in_bufz *C.size_t, out_bufsz *C.size_t) {
	fmt.Println("wbcuse_ioctl", req, cmd, arg, fi, flags, in_buf, in_bufz, out_bufsz)
}

//export wbcuse_poll
func wbcuse_poll(req *C.fuse_req_t, fi *C.fuse_file_info, ph *C.fuse_pollhandle) {
	fmt.Println("wbcuse_poll", req, fi)
}

func main() {
	args := []*C.char{}
	for _, arg := range os.Args {
		args = append(args, C.CString(arg))
	}

	panic(C.wbcuse_start(C.int(len(args)), &args[0]))
}
