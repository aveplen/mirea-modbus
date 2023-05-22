package main

import (
	"github.com/simonvetter/modbus"
)

type FallbackMiddleware struct {
	base   modbus.RequestHandler
	logger *Logger
}

func NewFallbackMiddleware(
	logger *Logger,
	base modbus.RequestHandler,
) *FallbackMiddleware {
	middleware := &FallbackMiddleware{
		logger: logger,
		base:   base,
	}

	return middleware
}

func (h *FallbackMiddleware) HandleCoils(req *modbus.CoilsRequest) ([]bool, error) {
	coils, err := h.base.HandleCoils(req)
	if coils != nil {
		return coils, err
	}

	h.logger.Debugf("HandleCoils returned 'nil' instead of coils, falling back to []bool{false}")
	return []bool{false}, err
}

func (h *FallbackMiddleware) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) ([]bool, error) {
	inputs, err := h.base.HandleDiscreteInputs(req)
	if inputs != nil {
		return inputs, err
	}

	h.logger.Debugf("HandleDiscreteInputs returned 'nil' instead of inputs, falling back to []bool{false}")
	return []bool{false}, err
}

func (h *FallbackMiddleware) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) ([]uint16, error) {
	registers, err := h.base.HandleHoldingRegisters(req)
	if registers != nil {
		return registers, err
	}

	h.logger.Debugf("HandleHoldingRegisters returned 'nil' instead of registers, falling back to []uint16{0}")
	return []uint16{0}, err
}

func (h *FallbackMiddleware) HandleInputRegisters(req *modbus.InputRegistersRequest) ([]uint16, error) {
	registers, err := h.base.HandleInputRegisters(req)
	if registers != nil {
		return registers, err
	}

	h.logger.Debugf("HandleInputRegisters returned 'nil' instead of registers, falling back to []uint16{0}")
	return []uint16{0}, err
}
