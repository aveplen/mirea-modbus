package main

type ClientManagmentService interface {
	ConnectParams(transport, address, port string) error
	Reconnect() error
	Disconnect() error
}

type ModbusService interface {
	ReadCoils0x01(addr uint16, cnt int) ([]bool, error)
	ReadDiscreteInputs0x02(addr uint16, cnt int) ([]bool, error)
	ReadHoldingRegisters0x03(addr uint16, cnt int) ([]uint16, error)
	ReadInputRegisters0x04(addr uint16, cnt int) ([]uint16, error)
	WriteSingleCoil0x05(addr uint16, value bool) error
	WriteSingleRegister0x06(addr uint16, value uint16) error
	WriteMultipleRegisters0x10(addr uint16, values []uint16) error
	WriteMultipleCoils0x0F(addr uint16, values []bool) error
}

type MainModelImpl struct {
	modbusService ModbusService
	clientService ClientManagmentService
}

func NewMainModelImpl(
	modbusService ModbusService,
	clientService ClientManagmentService,
) *MainModelImpl {
	return &MainModelImpl{
		modbusService: modbusService,
		clientService: clientService,
	}
}

func (m *MainModelImpl) Connect(transport, address, port string) error {
	return m.clientService.ConnectParams(transport, address, port)
}

func (m *MainModelImpl) Reconnect() error {
	return m.clientService.Reconnect()
}

func (m *MainModelImpl) Disconnect() error {
	return m.clientService.Disconnect()
}

func (m *MainModelImpl) ReadCoils(addr uint16, cnt int) ([]bool, error) {
	return m.modbusService.ReadCoils0x01(addr, cnt)
}

func (m *MainModelImpl) ReadDiscreteInputs(addr uint16, cnt int) ([]bool, error) {
	return m.modbusService.ReadDiscreteInputs0x02(addr, cnt)
}

func (m *MainModelImpl) ReadHoldingRegisters(addr uint16, cnt int) ([]uint16, error) {
	return m.modbusService.ReadHoldingRegisters0x03(addr, cnt)
}

func (m *MainModelImpl) ReadInputRegisters(addr uint16, cnt int) ([]uint16, error) {
	return m.modbusService.ReadInputRegisters0x04(addr, cnt)
}

func (m *MainModelImpl) WriteSingleCoil(addr uint16, value bool) error {
	return m.modbusService.WriteSingleCoil0x05(addr, value)
}

func (m *MainModelImpl) WriteSingleRegister(addr uint16, value uint16) error {
	return m.modbusService.WriteSingleRegister0x06(addr, value)
}

func (m *MainModelImpl) WriteMultipleRegisters(addr uint16, values []uint16) error {
	return m.modbusService.WriteMultipleRegisters0x10(addr, values)
}

func (m *MainModelImpl) WriteMultipleCoils(addr uint16, values []bool) error {
	return m.modbusService.WriteMultipleCoils0x0F(addr, values)
}
