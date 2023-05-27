package main

import "log"

type ServerManagerInterface interface {
	StartServer() error
	StopServer() error
}

type MainViewModel struct {
	serverManager ServerManagerInterface
}

func NewMainViewModel(serverManager ServerManagerInterface) *MainViewModel {
	return &MainViewModel{
		serverManager: serverManager,
	}
}

func (m *MainViewModel) StartServer() bool {
	if err := m.serverManager.StartServer(); err != nil {
		log.Printf("Could not start server, reason: %v", err)
		return false
	}

	log.Println("Server started successfuly")
	return true
}

func (m *MainViewModel) StopServer() bool {
	if err := m.serverManager.StopServer(); err != nil {
		log.Printf("Could not stop server, reason: %v", err)
		return false
	}

	log.Println("Server stopped successfuly")
	return true
}
