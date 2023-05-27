package main

import (
	"fmt"
	"log"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	seed, err := ReadSeed("seed.json")
	if err != nil {
		panic(fmt.Errorf("couls not read seed for modbus handler: %w", err))
	}

	service := NewModbusService(seed)
	fallback := NewFallbackMiddleware(
		NewValidationMiddleware(
			NewAdapterHandler(
				NewModbusHandler(service))))

	serverManager := NewServerManager(
		&modbus.ServerConfiguration{
			URL:        "tcp://localhost:5502",
			Timeout:    30 * time.Second,
			MaxClients: 5,
		},
		fallback,
	)

	activitySimulator := NewActivitySimulatorImpl(service, seed)

	viewModel := NewMainViewModel(serverManager, activitySimulator)
	view := NewView(seed, viewModel)

	service.SubscribeToCoilChanges(view.UpdateCoils)
	service.SubscribeToDiscreteInputChages(view.UpdateDiscreteInputs)
	service.SubscribeToHoldingRegisterChanges(view.UpdateHoldingRegisters)
	service.SubscribeToInputRegisterChanges(view.UpdateInputRegisters)
	log.SetOutput(&LogWriter{append: view.AppendLog})

	view.MainWindow.Run()
}
