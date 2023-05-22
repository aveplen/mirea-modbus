package main

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

type ModbusClient interface {
	Open() (err error)
	Close() (err error)

	ReadCoils(addr uint16, quantity uint16) (values []bool, err error)
	ReadCoil(addr uint16) (value bool, err error)

	ReadDiscreteInputs(addr uint16, quantity uint16) (values []bool, err error)
	ReadDiscreteInput(addr uint16) (value bool, err error)

	ReadRegisters(addr uint16, quantity uint16, regType modbus.RegType) (values []uint16, err error)
	ReadRegister(addr uint16, regType modbus.RegType) (value uint16, err error)

	WriteCoil(addr uint16, value bool) (err error)
	WriteCoils(addr uint16, values []bool) (err error)

	WriteRegister(addr uint16, value uint16) (err error)
	WriteRegisters(addr uint16, values []uint16) (err error)
}

type ModbusService struct {
	client ModbusClient
}

func NewModbusService(client ModbusClient) *ModbusService {
	return &ModbusService{
		client: client,
	}
}

func (a *ModbusService) SetClient(client ModbusClient) {
	a.client = client
}

func (a *ModbusService) ReadCoils0x01(addr uint16, cnt int) ([]bool, error) {
	fmt.Printf("ModbusService: client = %v\n", a.client)
	coils, err := a.client.ReadCoils(addr, uint16(cnt))
	if err != nil {
		return nil, fmt.Errorf("read %d coils at address %d: %w", cnt, addr, err)
	}
	return coils, nil
}

func (a *ModbusService) ReadDiscreteInputs0x02(addr uint16, cnt int) ([]bool, error) {
	inputs, err := a.client.ReadDiscreteInputs(addr, uint16(cnt))
	if err != nil {
		return nil, fmt.Errorf("read %d discrete inputs at address %d: %w", cnt, addr, err)
	}
	return inputs, nil
}

func (a *ModbusService) ReadHoldingRegisters0x03(addr uint16, cnt int) ([]uint16, error) {
	regs, err := a.client.ReadRegisters(addr, uint16(cnt), modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, fmt.Errorf("read %d holding registers at address %d: %w", cnt, addr, err)
	}
	return regs, nil
}

func (a *ModbusService) ReadInputRegisters0x04(addr uint16, cnt int) ([]uint16, error) {
	registers, err := a.client.ReadRegisters(addr, uint16(cnt), modbus.INPUT_REGISTER)
	if err != nil {
		return nil, fmt.Errorf("read %d input registers at address %d: %w", cnt, addr, err)
	}
	return registers, nil
}

func (a *ModbusService) WriteSingleCoil0x05(addr uint16, value bool) error {
	if err := a.client.WriteCoil(addr, value); err != nil {
		return fmt.Errorf("write coil at %d: %w", addr, err)
	}
	return nil
}

func (a *ModbusService) WriteSingleRegister0x06(addr uint16, value uint16) error {
	if err := a.client.WriteRegister(addr, value); err != nil {
		return fmt.Errorf("write value at %d: %w", addr, err)
	}
	return nil
}

func (a *ModbusService) WriteMultipleRegisters0x10(addr uint16, values []uint16) error {
	if err := a.client.WriteRegisters(addr, values); err != nil {
		return fmt.Errorf("write %d registers at address %d: %w", len(values), addr, err)
	}
	return nil
}

func (a *ModbusService) WriteMultipleCoils0x0F(addr uint16, values []bool) error {
	if err := a.client.WriteCoils(addr, values); err != nil {
		return fmt.Errorf("write %v coils at address %d: %w", len(values), addr, err)
	}
	return nil
}
