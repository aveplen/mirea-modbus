package main

import "log"

type ServerManagerInterface interface {
	StartServer() error
	StopServer() error
}

type ActivitySimulator interface {
	StartSimulation()
	StopSimulation()
}

type MainViewModel struct {
	serverManager     ServerManagerInterface
	activitySimulator ActivitySimulator
}

func NewMainViewModel(
	serverManager ServerManagerInterface,
	activitySimulator ActivitySimulator,
) *MainViewModel {
	return &MainViewModel{
		serverManager:     serverManager,
		activitySimulator: activitySimulator,
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

func (m *MainViewModel) StartSimulation() {
	m.activitySimulator.StartSimulation()
}

func (m *MainViewModel) StopSimulation() {
	m.activitySimulator.StopSimulation()
}
