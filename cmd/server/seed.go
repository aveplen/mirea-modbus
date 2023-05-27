package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

var (
	ErrNoInputCoils     = errors.New("no input coils object")
	ErrNoCoils          = errors.New("no coils object")
	ErrNoInputRegisters = errors.New("no input registers")
	ErrNoRegisters      = errors.New("no registers")

	ErrCoilsWrongType     = errors.New("coils object inconsistent typing")
	ErrRegistersWrongType = errors.New("registers object inconsistent typing")
)

func ReadSeed(filename string) (Dump, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Dump{}, fmt.Errorf("read file: %w", err)
	}

	seed := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &seed); err != nil {
		return Dump{}, fmt.Errorf("unmarshall seed file: %w", err)
	}

	coilsObj, ok := seed["coils"]
	if !ok {
		return Dump{}, ErrNoInputCoils
	}

	coils, err := unmarshallCoils(coilsObj)
	if err != nil {
		return Dump{}, err
	}

	discreteInputsObj, ok := seed["discrete_inputs"]
	if !ok {
		return Dump{}, ErrNoInputCoils
	}

	discreteInputs, err := unmarshallCoils(discreteInputsObj)
	if err != nil {
		return Dump{}, err
	}

	inputRegistersObj, ok := seed["input_registers"]
	if !ok {
		return Dump{}, ErrNoInputCoils
	}

	inpuRegisters, err := unmarhsallRegisters(inputRegistersObj)
	if err != nil {
		return Dump{}, err
	}

	holdingRegistersObj, ok := seed["holding_registers"]
	if !ok {
		return Dump{}, ErrNoInputCoils
	}

	holdingRegisters, err := unmarhsallRegisters(holdingRegistersObj)
	if err != nil {
		return Dump{}, err
	}

	return Dump{
		Coils:            coils,
		DiscreteInputs:   discreteInputs,
		InputRegisters:   inpuRegisters,
		HoldingRegisters: holdingRegisters,
	}, nil
}

func unmarshallCoils(intface interface{}) ([]Coil, error) {
	obj, ok := intface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't parse obj %v: %w", intface, ErrCoilsWrongType)
	}

	result := make([]Coil, 0, len(obj))
	for k, v := range obj {
		addr, err := strconv.Atoi(k)
		if err != nil {
			return nil, fmt.Errorf("atoi: %w", err)
		}

		coil, ok := v.(bool)
		if !ok {
			fmt.Println(k, v)
			return nil, fmt.Errorf("can't parse bool %v: %w", v, ErrCoilsWrongType)
		}

		result = append(result, Coil{
			addr:  uint16(addr),
			value: coil,
		})
	}

	return result, nil
}

func unmarhsallRegisters(intface interface{}) ([]Register, error) {
	obj, ok := intface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't parse obj %v: %w", intface, ErrRegistersWrongType)
	}

	result := make([]Register, 0, len(obj))
	for k, v := range obj {
		addr, err := strconv.Atoi(k)
		if err != nil {
			return nil, fmt.Errorf("atoi: %w", err)
		}

		register, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("can't parse float64 %v: %w", v, ErrRegistersWrongType)
		}

		result = append(result, Register{
			addr:  uint16(addr),
			value: uint16(register),
		})
	}

	return result, nil
}
