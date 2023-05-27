package main

import (
	"context"
	"math/rand"
	"time"
)

type ActivitySimulatorImpl struct {
	service *ModbusService
	seed    Dump

	cancel func()
}

func NewActivitySimulatorImpl(
	service *ModbusService,
	seed Dump,
) *ActivitySimulatorImpl {
	return &ActivitySimulatorImpl{
		service: service,
		seed:    seed,
	}
}

func (a *ActivitySimulatorImpl) StartSimulation() {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	go a.SimulateActivity(ctx)
}

func (a *ActivitySimulatorImpl) StopSimulation() {
	a.cancel()
}

func (a *ActivitySimulatorImpl) SimulateActivity(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			a.RandomizeDiscreteInputs()
			a.RandomizeHoldingRegisters()
		}
	}
}

func (a *ActivitySimulatorImpl) RandomizeDiscreteInputs() {
	for _, coil := range a.seed.Coils {
		a.service.SetCoil(coil.addr, rand.Int()%2 == 0)
	}
}

func (a *ActivitySimulatorImpl) RandomizeHoldingRegisters() {
	for _, coil := range a.seed.HoldingRegisters {
		a.service.SetHoldingRegister(coil.addr, uint16(rand.Int()))
	}
}
