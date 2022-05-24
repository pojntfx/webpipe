package cuse

import (
	"fmt"
	"sync"
)

var GlobalRegistry = &Registry{
	devices: map[int]Device{},
}

type Device interface {
	Init(registryID int, userdata Void, conn Conn)
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
