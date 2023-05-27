package main

import "github.com/simonvetter/modbus"

type ModbusService struct {
	coils            []Coil
	discreteInputs   []Coil
	holdingRegisters []Register
	inputRegisters   []Register

	coilSubs            []CoilSub
	discreteInputSubs   []CoilSub
	holdingRegisterSubs []RegisterSub
	inputRegisterSubs   []RegisterSub
}

type Dump struct {
	Coils            []Coil
	DiscreteInputs   []Coil
	HoldingRegisters []Register
	InputRegisters   []Register
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

func NewModbusService(seed Dump) *ModbusService {
	return &ModbusService{
		coils:            append(make([]Coil, 0, len(seed.Coils)), seed.Coils...),
		discreteInputs:   append(make([]Coil, 0, len(seed.DiscreteInputs)), seed.DiscreteInputs...),
		holdingRegisters: append(make([]Register, 0, len(seed.HoldingRegisters)), seed.HoldingRegisters...),
		inputRegisters:   append(make([]Register, 0, len(seed.InputRegisters)), seed.InputRegisters...),

		coilSubs:            make([]CoilSub, 0),
		discreteInputSubs:   make([]CoilSub, 0),
		holdingRegisterSubs: make([]RegisterSub, 0),
		inputRegisterSubs:   make([]RegisterSub, 0),
	}
}

func (s *ModbusService) GetCoil(addr uint16) (bool, error) {
	for _, v := range s.coils {
		if v.addr == addr {
			return v.value, nil
		}
	}
	return false, modbus.ErrIllegalDataAddress
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
		return modbus.ErrIllegalDataAddress
	}

	prev := s.coils[index].value
	s.coils[index].value = value
	for _, sub := range s.coilSubs {
		sub(CoilChange{from: prev, to: value, addr: addr})
	}

	return nil
}

func (s *ModbusService) GetDiscreteInputs(addr uint16) (bool, error) {
	for _, v := range s.discreteInputs {
		if v.addr == addr {
			return v.value, nil
		}
	}
	return false, modbus.ErrIllegalDataAddress
}

func (s *ModbusService) SetDiscreteInputs(addr uint16, value bool) error {
	var found bool
	var index int
	for i, v := range s.discreteInputs {
		if v.addr == addr {
			found = true
			index = i
			break
		}
	}

	if !found {
		return modbus.ErrIllegalDataAddress
	}

	prev := s.coils[index].value
	s.coils[index].value = value
	for _, sub := range s.coilSubs {
		sub(CoilChange{from: prev, to: value, addr: addr})
	}

	return nil
}

func (s *ModbusService) GetHoldingRegister(addr uint16) (uint16, error) {
	for _, v := range s.holdingRegisters {
		if v.addr == addr {
			return v.value, nil
		}
	}
	return 0, modbus.ErrIllegalDataAddress
}

func (s *ModbusService) SetHoldingRegister(addr uint16, value uint16) error {
	var found bool
	var index int
	for i, v := range s.holdingRegisters {
		if v.addr == addr {
			found = true
			index = i
			break
		}
	}

	if !found {
		return modbus.ErrIllegalDataAddress
	}

	prev := s.holdingRegisters[index].value
	s.holdingRegisters[index].value = value
	for _, sub := range s.holdingRegisterSubs {
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
	return 0, modbus.ErrIllegalDataAddress
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
		return modbus.ErrIllegalDataAddress
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

func (s *ModbusService) SubscribeToDiscreteInputChages(sub CoilSub) {
	s.discreteInputSubs = append(s.discreteInputSubs, sub)
}

func (s *ModbusService) SubscribeToHoldingRegisterChanges(sub RegisterSub) {
	s.holdingRegisterSubs = append(s.holdingRegisterSubs, sub)
}

func (s *ModbusService) SubscribeToInputRegisterChanges(sub RegisterSub) {
	s.inputRegisterSubs = append(s.inputRegisterSubs, sub)
}
