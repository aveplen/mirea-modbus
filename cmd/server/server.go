package main

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

type ServerManager struct {
	config  *modbus.ServerConfiguration
	handler modbus.RequestHandler

	server *modbus.ModbusServer
}

func NewServerManager(
	config *modbus.ServerConfiguration,
	handler modbus.RequestHandler,
) *ServerManager {

	return &ServerManager{
		config:  config,
		handler: handler,
	}
}

func (s *ServerManager) StartServer() error {
	server, err := modbus.NewServer(s.config, s.handler)
	if err != nil {
		return fmt.Errorf("create server: %w", err)
	}

	s.server = server
	if err := s.server.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}

func (s *ServerManager) StopServer() error {
	if err := s.server.Stop(); err != nil {
		return fmt.Errorf("stop server: %w", err)
	}

	s.server = nil
	return nil
}
