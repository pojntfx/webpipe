package cuse

import (
	"fmt"
	"sync"
)

var GlobalRegistry = &Registry{
	devices: map[int]Device{},
}

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

type Registry struct {
	devices map[int]Device
	lock    sync.Mutex
}

func (r *Registry) AddDevice(id int, device Device) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.devices[id] = device
}

func (r *Registry) RemoveDevice(id int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.devices, id)
}

func (r *Registry) GetDevice(id int) (Device, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	device, ok := r.devices[id]
	if !ok {
		return nil, fmt.Errorf("could not find device with id %v", id)
	}

	return device, nil
}
