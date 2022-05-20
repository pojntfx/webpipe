package manualcuse

/*
#cgo pkg-config: fuse3
#include <fuse3/cuse_lowlevel.h>
*/
import "C"

type Req C.fuse_req_t
