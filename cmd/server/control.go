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
		log.Printf("could not start server: %v", err)
		return false
	}

	log.Println("server started successfuly")
	return true
}

func (m *MainViewModel) StopServer() bool {
	if err := m.serverManager.StopServer(); err != nil {
		log.Printf("could not stop server: %v", err)
		return false
	}

	log.Println("server stopped successfuly")
	return true
}
