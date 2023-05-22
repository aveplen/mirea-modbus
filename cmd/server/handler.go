package main

import (
	"fmt"
)

type ModbusHandler struct {
	logger  *Logger
	service *ModbusService
}

func NewModbusHandler(logger *Logger, service *ModbusService) *ModbusHandler {
	return &ModbusHandler{
		logger:  logger,
		service: service,
	}
}

func (h *ModbusHandler) ReadCoils0x01(addr uint16, cnt int) ([]bool, error) {
	result := make([]bool, 0, cnt)

	for a := addr; a < addr+uint16(cnt); a++ {
		coil, err := h.service.GetCoil(a)
		if err != nil {
			return nil, fmt.Errorf("get coil at address %d: %w", a, err)
		}

		result = append(result, coil)
	}

	return result, nil
}

func (h *ModbusHandler) ReadDiscreteInputs0x02(addr uint16, cnt int) ([]bool, error) {
	result := make([]bool, 0, cnt)

	for a := addr; a < addr+uint16(cnt); a++ {
		reg, err := h.service.GetInputRegister(a)
		if err != nil {
			return nil, fmt.Errorf("get input register at address %d: %w", a, err)
		}

		result = append(result, reg > 0)
	}

	return result, nil
}

func (h *ModbusHandler) ReadHoldingRegisters0x03(addr uint16, cnt int) ([]uint16, error) {
	result := make([]uint16, 0, cnt)

	for a := addr; a < addr+uint16(cnt); a++ {
		reg, err := h.service.GetRegister(a)
		if err != nil {
			return nil, fmt.Errorf("get register at address %d: %w", a, err)
		}

		result = append(result, reg)
	}

	return result, nil
}

func (h *ModbusHandler) ReadInputRegisters0x04(addr uint16, cnt int) ([]uint16, error) {
	result := make([]uint16, 0, cnt)

	for a := addr; a < addr+uint16(cnt); a++ {
		inputReg, err := h.service.GetInputRegister(a)
		if err != nil {
			return nil, fmt.Errorf("get input register at address %d: %w", a, err)
		}

		result = append(result, inputReg)
	}

	return result, nil
}

func (h *ModbusHandler) WriteSingleCoil0x05(addr uint16, value bool) error {
	if err := h.service.SetCoil(addr, value); err != nil {
		return fmt.Errorf("set coil at address %d: %w", addr, err)
	}

	return nil
}

func (h *ModbusHandler) WriteSingleRegister0x06(addr uint16, value uint16) error {
	if err := h.service.SetRegister(addr, value); err != nil {
		return fmt.Errorf("set register at address %d: %w", addr, err)
	}

	return nil
}

func (h *ModbusHandler) WriteMultipleRegisters0x10(addr uint16, values []uint16) error {
	for i, value := range values {
		a := addr + uint16(i)

		if err := h.service.SetRegister(a, value); err != nil {
			return fmt.Errorf("set register at address %d: %w", a, err)
		}
	}

	return nil
}

func (h *ModbusHandler) WriteMultipleCoils0x0F(addr uint16, coils []bool) error {
	for i, coil := range coils {
		a := addr + uint16(i)

		if err := h.service.SetCoil(a, coil); err != nil {
			return fmt.Errorf("set coil at address %d: %w", a, err)
		}
	}

	return nil
}
