package main

import (
	"fmt"
	"log"

	"github.com/simonvetter/modbus"
)

type AdapterHandler struct {
	handler *ModbusHandler
}

func NewAdapterHandler(handler *ModbusHandler) *AdapterHandler {

	adapter := &AdapterHandler{
		handler: handler,
	}

	return adapter
}

// HandleCoils handles the read coils (0x01), write single coil (0x05)
// and write multiple coils (0x0F)
// - res:	coil values (only for reads)
// - err:	either nil if no error occurred, a modbus error
func (h *AdapterHandler) HandleCoils(req *modbus.CoilsRequest) ([]bool, error) {
	if req.IsWrite && req.Quantity == 1 {
		log.Printf(
			"Function 0x05 (write single coil) accessed with Addr: %d, Arguments: %v",
			req.Addr, req.Args,
		)

		if err := h.handler.WriteSingleCoil0x05(req.Addr, req.Args[0]); err != nil {
			return nil, fmt.Errorf("handle coils: %w", err)
		}

		return nil, nil
	}

	if req.IsWrite {
		log.Printf(
			"Function 0x0F (write multiple coils) accessed with Addr: %d, Quantity: %d, Arguments: %v",
			req.Addr, req.Quantity, req.Args,
		)

		if err := h.handler.WriteMultipleCoils0x0F(req.Addr, req.Args[:req.Quantity]); err != nil {
			return nil, fmt.Errorf("handle coils: %w", err)
		}

		return nil, nil
	}

	log.Printf(
		"Funtion 0x01 (read coils) accessed with Addr: %d, Quantity: %d",
		req.Addr, req.Quantity,
	)

	coils, err := h.handler.ReadCoils0x01(req.Addr, int(req.Quantity))
	if err != nil {
		return nil, fmt.Errorf("handle coils: %w", err)
	}

	return coils, nil
}

// HandleDiscreteInputs handles the read discrete inputs (0x02)
// - res: discrete input values
// - err:	either nil if no error occurred, a modbus error
func (h *AdapterHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) ([]bool, error) {
	inputs, err := h.handler.ReadDiscreteInputs0x02(req.Addr, int(req.Quantity))
	if err != nil {
		return nil, fmt.Errorf("handle discrete inputs: %w", err)
	}
	return inputs, nil
}

// HandleHoldingRegisters handles the read holding registers (0x03),
// write single register (0x06) and write multiple registers (0x10).
// A HoldingRegistersRequest object is passed to the handler (see above).
// - res:	register values
// - err:	either nil if no error occurred, a modbus error
func (h *AdapterHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) ([]uint16, error) {
	if req.IsWrite && req.Quantity == 1 {
		log.Printf(
			"Function 0x06 (write single register) accessed with Addr: %d, Quantity: %d, Arguments: %v",
			req.Addr, req.Quantity, req.Args,
		)

		if err := h.handler.WriteSingleRegister0x06(req.Addr, req.Args[0]); err != nil {
			return nil, fmt.Errorf("handle holding registers: %w", err)
		}

		return nil, nil
	}

	if req.IsWrite {
		log.Printf(
			"Function 0x10 (write mutliple registers) accessed with Addr: %d, Quantity: %d, Arguments: %v",
			req.Addr, req.Quantity, req.Args,
		)

		if err := h.handler.WriteMultipleRegisters0x10(req.Addr, req.Args[:req.Quantity]); err != nil {
			return nil, fmt.Errorf("handle holding registers: %w", err)
		}

		return nil, nil
	}

	log.Printf(
		"Function 0x03 (read holding registers) accessed with Addr: %d, Quantity: %d",
		req.Addr, req.Quantity,
	)

	regs, err := h.handler.ReadHoldingRegisters0x03(req.Addr, int(req.Quantity))
	if err != nil {
		return nil, fmt.Errorf("handle holding registers: %w", err)
	}

	return regs, nil
}

// HandleInputRegisters handles the read input registers (0x04)
// Note that input registers are always read-only as per the modbus spec.
// - res:	register values
// - err:	either nil if no error occurred, a modbus error
func (h *AdapterHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) ([]uint16, error) {
	log.Printf("Function 0x04 (read input registers) at Addr: %d with Quantity: %d", req.Addr, req.Quantity)
	regs, err := h.handler.ReadInputRegisters0x04(req.Addr, int(req.Quantity))
	if err != nil {
		return nil, fmt.Errorf("handle input registers: %w", err)
	}
	return regs, nil
}
