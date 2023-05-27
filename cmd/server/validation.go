package main

import (
	"log"

	"github.com/simonvetter/modbus"
)

type ValidationMiddleware struct {
	base modbus.RequestHandler
}

func NewValidationMiddleware(base modbus.RequestHandler) *ValidationMiddleware {
	middleware := &ValidationMiddleware{
		base: base,
	}
	return middleware
}

func (h *ValidationMiddleware) HandleCoils(req *modbus.CoilsRequest) ([]bool, error) {
	if req.UnitId != 1 {
		log.Printf("HandleCoils accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}
	return h.base.HandleCoils(req)
}

func (h *ValidationMiddleware) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) ([]bool, error) {
	if req.UnitId != 1 {
		log.Printf("HandleDiscreteInputs accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}
	return h.base.HandleDiscreteInputs(req)
}

func (h *ValidationMiddleware) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) ([]uint16, error) {
	if req.UnitId != 1 {
		log.Printf("HandleHoldingRegisters accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}
	return h.base.HandleHoldingRegisters(req)
}

func (h *ValidationMiddleware) HandleInputRegisters(req *modbus.InputRegistersRequest) ([]uint16, error) {
	if req.UnitId != 1 {
		log.Printf("HandleInputRegisters accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}
	return h.base.HandleInputRegisters(req)
}
