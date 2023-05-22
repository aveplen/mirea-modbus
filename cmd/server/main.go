package main

import (
	"fmt"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	seed := Dump{
		Coils: []Coil{
			{
				addr:  0x123,
				value: true,
			},
			{
				addr:  0x567,
				value: false,
			},
		},
	}

	logger := NewLogger(
		"SERVER",
		func() string { return time.Now().String() },
		100,
		10,
	)

	service := NewModbusService(logger, seed)
	handler := NewModbusHandler(logger, service)
	adapter := NewAdapterHandler(logger, handler)
	validation := NewValidationMiddleware(logger, adapter)
	fallback := NewFallbackMiddleware(logger, validation)

	server, err := modbus.NewServer(
		&modbus.ServerConfiguration{
			URL:        "tcp://localhost:5502",
			Timeout:    30 * time.Second,
			MaxClients: 5,
		},
		fallback,
	)

	if err != nil {
		panic(fmt.Errorf("failed to create server: %w", err))
	}

	startServer := func() bool {
		if err := server.Start(); err != nil {
			panic(fmt.Errorf("failed to start server: %w", err))
		}
		fmt.Println("server started")
		return true
	}

	stopServer := func() bool {
		if err := server.Stop(); err != nil {
			panic(fmt.Errorf("failed to stop server: %w", err))
		}
		fmt.Println("server stopped")
		return true
	}

	view := NewView(seed, startServer, stopServer)

	service.SubscribeToCoilChanges(func(change CoilChange) {
		fmt.Printf("coil change: %v\n", change)
		view.CoilsUpdateCallback(change)
	})

	service.SubscribeToRegisterChanges(func(change RegisterChange) {
		fmt.Printf("register change: %v\n", change)
		view.RegistersRenderCallback(change)
	})

	service.SubscribeToInputRegisterChanges(func(change RegisterChange) {
		fmt.Printf("input regisiter chnge: %v", change)
		view.InputRegistersRenderCallback(change)
	})

	view.MainWindow.Run()
}
