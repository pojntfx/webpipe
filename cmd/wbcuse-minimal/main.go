package main

/*
#cgo pkg-config: fuse3
#cgo CFLAGS: -Wno-error=implicit-function-declaration
#include "cuse.h"
*/
import (
	"C"
)
import (
	"fmt"
	"os"
)

//export wbcuse_open
func wbcuse_open(p0 *C.fuse_req_t, p1 *C.fuse_file_info) {
	fmt.Println("Hello, world!")

	return
}

func main() {
	args := []*C.char{}
	for _, arg := range os.Args {
		args = append(args, C.CString(arg))
	}

	panic(C.wbcuse_start(C.int(len(args)), &args[0]))
}
