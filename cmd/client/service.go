package main

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

type ClientSupplier interface {
	GetClient() (*modbus.ModbusClient, error)
}

type ModbusServiceImpl struct {
	clientService ClientSupplier
}

func NewModbusServiceImpl(clientService ClientSupplier) *ModbusServiceImpl {
	return &ModbusServiceImpl{
		clientService: clientService,
	}
}

func (a *ModbusServiceImpl) ReadCoils0x01(addr uint16, cnt int) ([]bool, error) {
	client, err := a.clientService.GetClient()
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	coils, err := client.ReadCoils(addr, uint16(cnt))
	if err != nil {
		return nil, fmt.Errorf("read %d coils at address %d: %w", cnt, addr, err)
	}
	return coils, nil
}

func (a *ModbusServiceImpl) ReadDiscreteInputs0x02(addr uint16, cnt int) ([]bool, error) {
	client, err := a.clientService.GetClient()
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	inputs, err := client.ReadDiscreteInputs(addr, uint16(cnt))
	if err != nil {
		return nil, fmt.Errorf("read %d discrete inputs at address %d: %w", cnt, addr, err)
	}
	return inputs, nil
}

func (a *ModbusServiceImpl) ReadHoldingRegisters0x03(addr uint16, cnt int) ([]uint16, error) {
	client, err := a.clientService.GetClient()
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	regs, err := client.ReadRegisters(addr, uint16(cnt), modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, fmt.Errorf("read %d holding registers at address %d: %w", cnt, addr, err)
	}
	return regs, nil
}

func (a *ModbusServiceImpl) ReadInputRegisters0x04(addr uint16, cnt int) ([]uint16, error) {
	client, err := a.clientService.GetClient()
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	registers, err := client.ReadRegisters(addr, uint16(cnt), modbus.INPUT_REGISTER)
	if err != nil {
		return nil, fmt.Errorf("read %d input registers at address %d: %w", cnt, addr, err)
	}
	return registers, nil
}

func (a *ModbusServiceImpl) WriteSingleCoil0x05(addr uint16, value bool) error {
	client, err := a.clientService.GetClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	if err := client.WriteCoil(addr, value); err != nil {
		return fmt.Errorf("write coil at %d: %w", addr, err)
	}
	return nil
}

func (a *ModbusServiceImpl) WriteSingleRegister0x06(addr uint16, value uint16) error {
	client, err := a.clientService.GetClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	if err := client.WriteRegister(addr, value); err != nil {
		return fmt.Errorf("write value at %d: %w", addr, err)
	}
	return nil
}

func (a *ModbusServiceImpl) WriteMultipleRegisters0x10(addr uint16, values []uint16) error {
	client, err := a.clientService.GetClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	if err := client.WriteRegisters(addr, values); err != nil {
		return fmt.Errorf("write %d registers at address %d: %w", len(values), addr, err)
	}
	return nil
}

func (a *ModbusServiceImpl) WriteMultipleCoils0x0F(addr uint16, values []bool) error {
	client, err := a.clientService.GetClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	if err := client.WriteCoils(addr, values); err != nil {
		return fmt.Errorf("write %v coils at address %d: %w", len(values), addr, err)
	}
	return nil
}
