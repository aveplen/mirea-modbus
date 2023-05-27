package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/simonvetter/modbus"
)

const (
	DefaultTransport = "tcp"
	DefaultAddress   = "localhost"
	DefaultPort      = "5502"
)

var (
	ErrTransportUnknown = errors.New("transport unknown")
	ErrAddressUnknown   = errors.New("address unknown")
	ErrPortUnknown      = errors.New("port unknown")
	ErrNoClient         = errors.New("no client")
	ErrNotEstablished   = errors.New("connection not established")
)

type ClientManagmentServiceImpl struct {
	client *modbus.ModbusClient

	transport       string
	transportSet    bool
	address         string
	addressSet      bool
	port            string
	portSet         bool
	connEstablished bool

	logPrefix string
}

func NewClientManagmentSercieImpl() *ClientManagmentServiceImpl {
	return &ClientManagmentServiceImpl{}
}

func (m *ClientManagmentServiceImpl) SetParams(transport, address, port string) {
	m.SetTransport(transport)
	m.SetAddress(address)
	m.SetPort(port)
}

func (m *ClientManagmentServiceImpl) SetTransport(transport string) {
	m.transport = transport
	m.transportSet = true
}

func (m *ClientManagmentServiceImpl) SetAddress(address string) {
	m.address = address
	m.addressSet = true
}

func (m *ClientManagmentServiceImpl) SetPort(port string) {
	m.port = port
	m.portSet = true
}

func (m *ClientManagmentServiceImpl) ConnectParams(transport, address, port string) error {
	m.SetParams(transport, address, port)
	if err := m.Connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	return nil
}

func (m *ClientManagmentServiceImpl) ConnectDefalut() error {
	transport := m.resolveTransportDefault()
	address := m.resolveAddressDefault()
	port := m.resolvePortDefault()
	if err := m.ConnectParams(transport, address, port); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	return nil
}

func (m *ClientManagmentServiceImpl) Connect() error {
	m.client = nil
	m.connEstablished = false

	transport, err := m.resolveTransportStrict()
	if err != nil {
		return fmt.Errorf("resolve transport: %w", err)
	}

	address, err := m.resolveAddressStrict()
	if err != nil {
		return fmt.Errorf("resolve address: %w", err)
	}

	port, err := m.resolvePortStrict()
	if err != nil {
		return fmt.Errorf("resolve port: %w", err)
	}

	client, err := modbus.NewClient(
		&modbus.ClientConfiguration{
			URL: fmt.Sprintf("%s://%s:%s", transport, address, port),
		},
	)

	if err != nil {
		log.Printf("%s: client not created: %v\n", m.logPrefix, err)
		return err
	}

	if err := client.Open(); err != nil {
		log.Printf("%s: conn is not established: %v\n", m.logPrefix, err)
		return err
	}

	log.Printf("%s: conn established", m.logPrefix)
	m.connEstablished = true
	m.client = client
	return nil
}

func (m *ClientManagmentServiceImpl) Reconnect() error {
	if err := m.Disconnect(); err != nil {
		return fmt.Errorf("disconnect: %w", err)
	}

	if err := m.Connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	return nil
}

func (m *ClientManagmentServiceImpl) Disconnect() error {
	if m.connEstablished && m.client != nil {
		if err := m.client.Close(); err != nil {
			log.Printf("%s: close client connection: %v", m.logPrefix, err)
		}

		m.client = nil
		m.connEstablished = false
	}

	return nil
}

func (m *ClientManagmentServiceImpl) GetClient() (*modbus.ModbusClient, error) {
	if m.client == nil {
		return nil, ErrNoClient
	}

	if !m.connEstablished {
		return nil, ErrNotEstablished
	}

	return m.client, nil
}

func (m *ClientManagmentServiceImpl) resolveTransportDefault() string {
	if !m.transportSet {
		return DefaultTransport
	}
	return m.transport
}

func (m *ClientManagmentServiceImpl) resolveTransportStrict() (string, error) {
	if !m.transportSet {
		return "", ErrTransportUnknown
	}
	return m.transport, nil
}

func (m *ClientManagmentServiceImpl) resolveAddressDefault() string {
	if !m.addressSet {
		return DefaultAddress
	}
	return m.address
}

func (m *ClientManagmentServiceImpl) resolveAddressStrict() (string, error) {
	if !m.addressSet {
		return "", ErrAddressUnknown
	}
	return m.address, nil
}

func (m *ClientManagmentServiceImpl) resolvePortDefault() string {
	if !m.portSet {
		return DefaultPort
	}
	return m.port
}

func (m *ClientManagmentServiceImpl) resolvePortStrict() (string, error) {
	if !m.portSet {
		return "", ErrPortUnknown
	}
	return m.port, nil
}
