package main

import (
	"fmt"
	"log"

	"github.com/simonvetter/modbus"
)

const retries = 2

type ClientSupplier interface {
	GetClient() (*modbus.ModbusClient, error)
	Reconnect() error
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
	var coils []bool
	var err error
	var succ bool

	for attempt := 0; attempt < retries+1; attempt++ {
		coils, err = a.readCoils0x01(addr, cnt)
		if err == nil {
			succ = true
			break
		}

		log.Printf(
			"retry reading %d coils at 0x%X, attempts left: %d",
			cnt, addr, retries-attempt,
		)
		a.clientService.Reconnect()
	}

	if !succ {
		return nil, err
	}

	return coils, nil
}

func (a *ModbusServiceImpl) readCoils0x01(addr uint16, cnt int) ([]bool, error) {
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
	var discreteInputs []bool
	var err error
	var succ bool

	for attempt := 0; attempt < retries+1; attempt++ {
		discreteInputs, err = a.readDiscreteInputs0x02(addr, cnt)
		if err == nil {
			succ = true
			break
		}

		log.Printf(
			"retry reading discrete inputs %d coils at 0x%X, attempts left: %d",
			cnt, addr, retries-attempt,
		)
		a.clientService.Reconnect()
	}

	if !succ {
		return nil, err
	}

	return discreteInputs, nil
}

func (a *ModbusServiceImpl) readDiscreteInputs0x02(addr uint16, cnt int) ([]bool, error) {
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
	var holdingRegisters []uint16
	var err error
	var succ bool

	for attempt := 0; attempt < retries+1; attempt++ {
		holdingRegisters, err = a.readHoldingRegisters0x03(addr, cnt)
		if err == nil {
			succ = true
			break
		}

		log.Printf(
			"retry reading holding registers %d coils at 0x%X, attempts left: %d",
			cnt, addr, retries-attempt,
		)
		a.clientService.Reconnect()
	}

	if !succ {
		return nil, err
	}

	return holdingRegisters, nil
}

func (a *ModbusServiceImpl) readHoldingRegisters0x03(addr uint16, cnt int) ([]uint16, error) {
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
	var inputRegisters []uint16
	var err error
	var succ bool

	for attempt := 0; attempt < retries+1; attempt++ {
		inputRegisters, err = a.readInputRegisters0x04(addr, cnt)
		if err == nil {
			succ = true
			break
		}

		log.Printf(
			"retry reading input registers %d coils at 0x%X, attempts left: %d",
			cnt, addr, retries-attempt,
		)
		a.clientService.Reconnect()
	}

	if !succ {
		return nil, err
	}

	return inputRegisters, nil
}

func (a *ModbusServiceImpl) readInputRegisters0x04(addr uint16, cnt int) ([]uint16, error) {
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
	var err error
	var succ bool

	for attempt := 0; attempt < retries+1; attempt++ {
		if err = a.writeSingleCoil0x05(addr, value); err == nil {
			succ = true
			break
		}

		log.Printf(
			"retry writing %v to single coil at 0x%X, attempts left: %d",
			value, addr, retries-attempt,
		)
		a.clientService.Reconnect()
	}

	if !succ {
		return err
	}

	return nil
}

func (a *ModbusServiceImpl) writeSingleCoil0x05(addr uint16, value bool) error {
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
	var err error
	var succ bool

	for attempt := 0; attempt < retries+1; attempt++ {
		if err = a.writeSingleRegister0x06(addr, value); err == nil {
			succ = true
			break
		}

		log.Printf(
			"retry writing %v to single register at 0x%X, attempts left: %d",
			value, addr, retries-attempt,
		)
		a.clientService.Reconnect()
	}

	if !succ {
		return err
	}

	return nil
}

func (a *ModbusServiceImpl) writeSingleRegister0x06(addr uint16, value uint16) error {
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
	var err error
	var succ bool

	for attempt := 0; attempt < retries+1; attempt++ {
		if err = a.writeMultipleRegisters0x10(addr, values); err == nil {
			succ = true
			break
		}

		log.Printf(
			"retry writing %v to multiple registers starting at 0x%X, attempts left: %d",
			values, addr, retries-attempt,
		)
		a.clientService.Reconnect()
	}

	if !succ {
		return err
	}

	return nil
}

func (a *ModbusServiceImpl) writeMultipleRegisters0x10(addr uint16, values []uint16) error {
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
	var err error
	var succ bool

	for attempt := 0; attempt < retries+1; attempt++ {
		if err = a.writeMultipleCoils0x0F(addr, values); err == nil {
			succ = true
			break
		}

		log.Printf(
			"retry writing %v to multiple coils starting at 0x%X, attempts left: %d",
			values, addr, retries-attempt,
		)
		a.clientService.Reconnect()
	}

	if !succ {
		return err
	}

	return nil
}

func (a *ModbusServiceImpl) writeMultipleCoils0x0F(addr uint16, values []bool) error {
	client, err := a.clientService.GetClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	if err := client.WriteCoils(addr, values); err != nil {
		return fmt.Errorf("write %v coils at address %d: %w", len(values), addr, err)
	}

	return nil
}
