package main

import "errors"

var (
	ErrNoSuchCoil          = errors.New("no such coil")
	ErrNoSuchRegister      = errors.New("no such register")
	ErrNoSuchInputRegister = errors.New("no such input register")
)

type ModbusService struct {
	logger *Logger

	coils          []Coil
	registers      []Register
	inputRegisters []Register

	coilSubs          []CoilSub
	registerSubs      []RegisterSub
	inputRegisterSubs []RegisterSub
}

type Dump struct {
	Coils          []Coil
	Registers      []Register
	InputRegisters []Register
}

type Coil struct {
	addr  uint16
	value bool
}

type CoilChange struct {
	addr uint16
	from bool
	to   bool
}

type CoilSub func(change CoilChange)

type Register struct {
	addr  uint16
	value uint16
}

type RegisterChange struct {
	addr uint16
	from uint16
	to   uint16
}

type RegisterSub func(change RegisterChange)

func NewModbusService(logger *Logger, seed Dump) *ModbusService {
	return &ModbusService{
		logger: logger,

		coils:          append(make([]Coil, 0, len(seed.Coils)), seed.Coils...),
		registers:      append(make([]Register, 0, len(seed.Registers)), seed.Registers...),
		inputRegisters: append(make([]Register, 0, len(seed.InputRegisters)), seed.InputRegisters...),

		coilSubs:          make([]CoilSub, 0),
		registerSubs:      make([]RegisterSub, 0),
		inputRegisterSubs: make([]RegisterSub, 0),
	}
}

func (s *ModbusService) GetCoil(addr uint16) (bool, error) {
	for _, v := range s.coils {
		if v.addr == addr {
			return v.value, nil
		}
	}
	return false, ErrNoSuchCoil
}

func (s *ModbusService) SetCoil(addr uint16, value bool) error {
	var found bool
	var index int
	for i, v := range s.coils {
		if v.addr == addr {
			found = true
			index = i
			break
		}
	}

	if !found {
		return ErrNoSuchCoil
	}

	prev := s.coils[index].value
	s.coils[index].value = value
	for _, sub := range s.coilSubs {
		sub(CoilChange{from: prev, to: value, addr: addr})
	}

	return nil
}

func (s *ModbusService) GetRegister(addr uint16) (uint16, error) {
	for _, v := range s.registers {
		if v.addr == addr {
			return v.value, nil
		}
	}
	return 0, ErrNoSuchCoil
}

func (s *ModbusService) SetRegister(addr uint16, value uint16) error {
	var found bool
	var index int
	for i, v := range s.registers {
		if v.addr == addr {
			found = true
			index = i
			break
		}
	}

	if !found {
		return ErrNoSuchRegister
	}

	prev := s.registers[index].value
	s.registers[index].value = value
	for _, sub := range s.registerSubs {
		sub(RegisterChange{from: prev, to: value, addr: addr})
	}

	return nil
}

func (s *ModbusService) GetInputRegister(addr uint16) (uint16, error) {
	for _, v := range s.inputRegisters {
		if v.addr == addr {
			return v.value, nil
		}
	}
	return 0, ErrNoSuchCoil
}

func (s *ModbusService) SetInputRegister(addr uint16, value uint16) error {
	var found bool
	var index int
	for i, v := range s.inputRegisters {
		if v.addr == addr {
			found = true
			index = i
			break
		}
	}

	if !found {
		return ErrNoSuchInputRegister
	}

	prev := s.inputRegisters[index].value
	s.inputRegisters[index].value = value
	for _, sub := range s.inputRegisterSubs {
		sub(RegisterChange{from: prev, to: value, addr: addr})
	}

	return nil
}

func (s *ModbusService) SubscribeToCoilChanges(sub CoilSub) {
	s.coilSubs = append(s.coilSubs, sub)
}

func (s *ModbusService) SubscribeToRegisterChanges(sub RegisterSub) {
	s.registerSubs = append(s.registerSubs, sub)
}

func (s *ModbusService) SubscribeToInputRegisterChanges(sub RegisterSub) {
	s.inputRegisterSubs = append(s.inputRegisterSubs, sub)
}
