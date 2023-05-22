package main

import (
	"github.com/simonvetter/modbus"
)

type ValidationMiddleware struct {
	base   modbus.RequestHandler
	logger *Logger
}

func NewValidationMiddleware(
	logger *Logger,
	base modbus.RequestHandler,
) *ValidationMiddleware {
	middleware := &ValidationMiddleware{
		logger: logger,
		base:   base,
	}
	return middleware
}

func (h *ValidationMiddleware) HandleCoils(req *modbus.CoilsRequest) ([]bool, error) {
	if req.UnitId != 1 {
		h.logger.Errorf("HandleCoils accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}
	return h.base.HandleCoils(req)
}

func (h *ValidationMiddleware) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) ([]bool, error) {
	if req.UnitId != 1 {
		h.logger.Errorf("HandleDiscreteInputs accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}
	return h.base.HandleDiscreteInputs(req)
}

func (h *ValidationMiddleware) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) ([]uint16, error) {
	if req.UnitId != 1 {
		h.logger.Errorf("HandleHoldingRegisters accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}
	return h.base.HandleHoldingRegisters(req)
}

func (h *ValidationMiddleware) HandleInputRegisters(req *modbus.InputRegistersRequest) ([]uint16, error) {
	if req.UnitId != 1 {
		h.logger.Errorf("HandleInputRegisters accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}
	return h.base.HandleInputRegisters(req)
}
