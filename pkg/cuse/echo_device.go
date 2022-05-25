package cuse

import "fmt"

type EchoDevice struct {
	Device

	registryID int
}

func NewEchoDevice(registryID int) *EchoDevice {
	return &EchoDevice{
		registryID: registryID,
	}
}

func (d *EchoDevice) OpenDevice(args []string) error {
	return StartCUSE(d.registryID, args)
}

func (d *EchoDevice) Init(userdata Void, conn Conn) {
	fmt.Println("Init", userdata, conn)
}

func (d *EchoDevice) InitDone(userdata Void) {
	fmt.Println("InitDone", userdata)
}

func (d *EchoDevice) Destroy(userdata Void) {
	fmt.Println("Destroy", userdata)
}

func (d *EchoDevice) Open(req Request, fi FileInfo) {
	fmt.Println("Open", req, fi)
}

func (d *EchoDevice) Read(req Request, size Size, off Offset, fi FileInfo) {
	fmt.Println("Read", req, size, off, fi)
}

func (d *EchoDevice) Write(req Request, buf Buffer, size Size, off Offset, fi FileInfo) {
	fmt.Println("Write", req, buf, size, off, fi)
}

func (d *EchoDevice) Flush(req Request, fi FileInfo) {
	fmt.Println("Flush", req, fi)
}

func (d *EchoDevice) Fsync(req Request, datasync int, fi FileInfo) {
	fmt.Println("Fsync", req, datasync, fi)
}

func (d *EchoDevice) Ioctl(req Request, cmd int, arg Void, fi FileInfo, flags uint, inputBuf Void, inputBufSize Size, outBufSize Size) {
	fmt.Println("Ioctl", req, cmd, arg, fi, flags, inputBuf, inputBufSize, outBufSize)
}

func (d *EchoDevice) Poll(req Request, fi FileInfo, ph PollHandle) {
	fmt.Println("Poll", req, fi, ph)
}
