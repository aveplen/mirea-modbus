package main

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

type Application struct {
	client  ModbusClient
	service *ModbusService

	transport       *string
	address         *string
	port            *string
	connEstablished bool
}

func (a *Application) SetClient(client ModbusClient) {
	a.service.SetClient(client)
	a.client = client
}

func (a *Application) Connect(transport, address, port string) error {
	a.transport = &transport
	a.address = &address
	a.port = &port

	client, err := modbus.NewClient(
		&modbus.ClientConfiguration{
			URL: fmt.Sprintf("%s://%s:%s", *a.transport, *a.address, *a.port),
		},
	)

	if err != nil {
		fmt.Printf("connect: client not created: %v\n", err)
		return err
	}

	a.SetClient(client)

	if err := a.client.Open(); err != nil {
		fmt.Printf("connect: conn is not established: %v\n", err)
		return err
	}

	fmt.Println("connect: conn established")
	a.connEstablished = true
	return nil
}

func (a *Application) Reconnect() error {
	if a.connEstablished {
		if err := a.client.Close(); err != nil {
			return err
		} else {
			fmt.Println("reconnect: conn is not established, skip closing")
		}
	}

	client, err := modbus.NewClient(
		&modbus.ClientConfiguration{
			URL: fmt.Sprintf("%s://%s:%s", *a.transport, *a.address, *a.port),
		},
	)

	if err != nil {
		fmt.Printf("reconnect: client not created: %v\n", err)
		return err
	}

	a.SetClient(client)

	if err := client.Open(); err != nil {
		fmt.Printf("reconnect: conn is not established: %v\n", err)
		return err
	}

	fmt.Println("reconnect: conn established")
	return nil
}

func (a *Application) Disconnect() error {
	if a.connEstablished {
		if err := a.client.Close(); err != nil {
			return err
		} else {
			fmt.Println("disconnect: conn is not established, skip closing")
		}
	}

	fmt.Println("disconnect: connection closed")
	return nil
}

func (a *Application) ReadCoils(addr uint16, cnt int) ([]bool, error) {
	return a.service.ReadCoils0x01(addr, cnt)
}

func (a *Application) ReadDiscreteInputs(addr uint16, cnt int) ([]bool, error) {
	return a.service.ReadDiscreteInputs0x02(addr, cnt)
}

func (a *Application) ReadHoldingRegisters(addr uint16, cnt int) ([]uint16, error) {
	return a.service.ReadHoldingRegisters0x03(addr, cnt)
}

func (a *Application) ReadInputRegisters(addr uint16, cnt int) ([]uint16, error) {
	return a.service.ReadInputRegisters0x04(addr, cnt)
}

func (a *Application) WriteSingleCoil(addr uint16, value bool) error {
	return a.service.WriteSingleCoil0x05(addr, value)
}

func (a *Application) WriteSingleRegister(addr uint16, value uint16) error {
	return a.service.WriteSingleRegister0x06(addr, value)
}

func (a *Application) WriteMultipleRegisters(addr uint16, values []uint16) error {
	return a.service.WriteMultipleRegisters0x10(addr, values)
}

func (a *Application) WriteMultipleCoils(addr uint16, values []bool) error {
	return a.service.WriteMultipleCoils0x0F(addr, values)
}

func (a *Application) AppMain() {
	a.service = NewModbusService(nil)

	NewView(ViewCallbacks{
		connect:                a.Connect,
		reconnect:              a.Reconnect,
		disconnect:             a.Disconnect,
		readCoils:              a.ReadCoils,
		readDiscreteInputs:     a.ReadDiscreteInputs,
		readHoldingRegisters:   a.ReadHoldingRegisters,
		readInputRegisters:     a.ReadInputRegisters,
		writeSingleCoil:        a.WriteSingleCoil,
		writeSingleRegister:    a.WriteSingleRegister,
		writeMultipleRegisters: a.WriteMultipleRegisters,
		writeMultipleCoils:     a.WriteMultipleCoils,
	}).Render().Run()
}

func main() {
	app := Application{}
	app.AppMain()
}
