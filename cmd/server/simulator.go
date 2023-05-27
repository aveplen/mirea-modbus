package main

import (
	"context"
	"log"
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
			a.RandomizeInputRegisters()
		}
	}
}

func (a *ActivitySimulatorImpl) RandomizeDiscreteInputs() {
	for _, coil := range a.seed.DiscreteInputs {
		c := rand.Intn(100)%2 == 0
		a.service.SetDiscreteInput(coil.addr, c)
		log.Printf("upating discret input at 0x%X to %v", coil.addr, c)
	}
}

func (a *ActivitySimulatorImpl) RandomizeInputRegisters() {
	for _, reg := range a.seed.InputRegisters {
		r := uint16(rand.Int())
		a.service.SetInputRegister(reg.addr, r)
		log.Printf("upating input reg at 0x%X to 0x%X", reg.addr, r)
	}
}
