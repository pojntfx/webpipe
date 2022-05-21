package main

/*
#cgo pkg-config: fuse3
#include "cuse.h"
*/
import (
	"C"
)
import (
	"os"
)

func main() {
	args := []*C.char{}
	for _, arg := range os.Args {
		args = append(args, C.CString(arg))
	}

	panic(C.wbcuse_start(C.int(len(args)), &args[0]))
}
