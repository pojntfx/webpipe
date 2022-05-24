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

func (d *EchoDevice) Open(args []string) error {
	return StartCUSE(d.registryID, args)
}

func (d *EchoDevice) Init(registryID int, userdata Void, conn Conn) {
	fmt.Println("Init", registryID, userdata, conn)
}
