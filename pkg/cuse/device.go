package cuse

type Device interface {
	Init(userdata Void, conn Conn)
	InitDone(userdata Void)
	Destroy(userdata Void)
	Open(req Request, fi FileInfo)
	Read(req Request, size Size, off Offset, fi FileInfo)
	Write(req Request, buf Buffer, size Size, off Offset, fi FileInfo)
	Flush(req Request, fi FileInfo)
	Release(req Request, fi FileInfo)
	Fsync(req Request, datasync int, fi FileInfo)
	Ioctl(req Request, cmd int, arg Void, fi FileInfo, flags uint, inputBuf Void, inputBufSize Size, outBufSize Size)
	Poll(req Request, fi FileInfo, ph PollHandle)
}

type deviceContainer struct {
	device Device
}
