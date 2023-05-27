package main

import (
	"fmt"
	"log"
)

type ModbusHandler struct {
	service *ModbusService
}

func NewModbusHandler(service *ModbusService) *ModbusHandler {
	return &ModbusHandler{
		service: service,
	}
}

func (h *ModbusHandler) ReadCoils0x01(addr uint16, cnt int) ([]bool, error) {
	log.Printf("Call function 0x01 (read coils), addr: 0x%X, cnt: %d", addr, cnt)

	result := make([]bool, 0, cnt)
	for a := addr; a < addr+uint16(cnt); a++ {
		coil, err := h.service.GetCoil(a)
		if err != nil {
			log.Printf("Could not get coil at addr: 0x%X, reason: %v", addr, err)
			return nil, fmt.Errorf("get coil at addr 0x%X: %w", a, err)
		}

		result = append(result, coil)
	}

	log.Printf("Successfuly read %d coils at addr: 0x%X", cnt, addr)
	return result, nil
}

func (h *ModbusHandler) ReadDiscreteInputs0x02(addr uint16, cnt int) ([]bool, error) {
	log.Printf("Call function 0x02 (read discrete inputs), addr: 0x%X, cnt: %d", addr, cnt)
	log.Printf("Discrete inputs: %v", h.service.discreteInputs)

	result := make([]bool, 0, cnt)
	for a := addr; a < addr+uint16(cnt); a++ {
		dis, err := h.service.GetDiscreteInputs(a)
		if err != nil {
			log.Printf("Could not get discrete input at addr: 0x%X, reason: %v", addr, err)
			return nil, fmt.Errorf("get discrete input at addr 0x%X: %w", a, err)
		}

		result = append(result, dis)
	}

	log.Printf("Successfuly read %d discrete inputs at addr: 0x%X", cnt, addr)
	return result, nil
}

func (h *ModbusHandler) ReadHoldingRegisters0x03(addr uint16, cnt int) ([]uint16, error) {
	log.Printf("Call function 0x03 (read holding registers), addr: 0x%X, cnt: %d", addr, cnt)

	result := make([]uint16, 0, cnt)
	for a := addr; a < addr+uint16(cnt); a++ {
		reg, err := h.service.GetHoldingRegister(a)
		if err != nil {
			log.Printf("Could not get holding register at addr: 0x%X, reason: %v", addr, err)
			return nil, fmt.Errorf("get register at addr 0x%X: %w", a, err)
		}

		result = append(result, reg)
	}

	log.Printf("Successfuly read %d holding registers at addr: 0x%X", cnt, addr)
	return result, nil
}

func (h *ModbusHandler) ReadInputRegisters0x04(addr uint16, cnt int) ([]uint16, error) {
	log.Printf("Call function 0x04 (read input registers), addr: 0x%X, cnt: %d", addr, cnt)

	result := make([]uint16, 0, cnt)
	for a := addr; a < addr+uint16(cnt); a++ {
		inputReg, err := h.service.GetInputRegister(a)
		if err != nil {
			log.Printf("Could not get input register at addr: 0x%X, reason: %v", addr, err)
			return nil, fmt.Errorf("get input register at addr 0x%X: %w", a, err)
		}

		result = append(result, inputReg)
	}

	log.Printf("Successfuly read %d input registers at addr: 0x%X", cnt, addr)
	return result, nil
}

func (h *ModbusHandler) WriteSingleCoil0x05(addr uint16, value bool) error {
	log.Printf("Call function 0x05 (write single coil), addr: 0x%X, value: %v", addr, value)

	if err := h.service.SetCoil(addr, value); err != nil {
		log.Printf("Could not write coil at addr: 0x%X, reason: %v", addr, err)
		return fmt.Errorf("set coil at addr 0x%X: %w", addr, err)
	}

	log.Printf("Successfuly written %v to coil at addr: 0x%X", value, addr)
	return nil
}

func (h *ModbusHandler) WriteSingleRegister0x06(addr uint16, value uint16) error {
	log.Printf("Call function 0x06 (write single register), addr: 0x%X, value: %d", addr, value)

	if err := h.service.SetHoldingRegister(addr, value); err != nil {
		log.Printf("Could not write register at addr: 0x%X, reason: %v", addr, err)
		return fmt.Errorf("set register at addr 0x%X: %w", addr, err)
	}

	log.Printf("Successfuly written %v to register at addr: 0x%X", value, addr)
	return nil
}

func (h *ModbusHandler) WriteMultipleRegisters0x10(addr uint16, values []uint16) error {
	log.Printf("Call function 0x10 (write multiple registers), addr: 0x%X, values: %v", addr, values)

	for i, value := range values {
		a := addr + uint16(i)

		if err := h.service.SetHoldingRegister(a, value); err != nil {
			log.Printf("Could not write multple registers at addr: 0x%X, reason: %v", addr, err)
			return fmt.Errorf("set register at addr 0x%X: %w", a, err)
		}
	}

	log.Printf("Successfuly written %v to registers at addr: 0x%X", values, addr)
	return nil
}

func (h *ModbusHandler) WriteMultipleCoils0x0F(addr uint16, coils []bool) error {
	log.Printf("Call function 0x0F (write multiple coils), addr: 0x%X, values: %v", addr, coils)

	for i, coil := range coils {
		a := addr + uint16(i)

		if err := h.service.SetCoil(a, coil); err != nil {
			log.Printf("Could not write coils at addr: 0x%X, reason: %v", addr, err)
			return fmt.Errorf("set coil at addr 0x%X: %w", a, err)
		}
	}

	log.Printf("Successfuly written %v to coils at addr: 0x%X", coils, addr)
	return nil
}
