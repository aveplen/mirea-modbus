package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/simonvetter/modbus"
)

type Application struct {
	client *modbus.ModbusClient
}

// 0x01 Read coils
func (a *Application) ReadCoils() ([]bool, error) {
	addr, err := askUint16("Starting address")
	if err != nil {
		return nil, fmt.Errorf("ask address: %w", err)
	}

	quantity, err := askUint16("Quantity")
	if err != nil {
		return nil, fmt.Errorf("ask quantity: %w", err)
	}

	coils, err := a.client.ReadCoils(addr, quantity)
	if err != nil {
		return nil, fmt.Errorf("could not read %d coils at %d: %w", quantity, addr, err)
	}

	return coils, nil
}

// 0x02 Read discrete inputs
func (a *Application) ReadDiscreteInputs() ([]bool, error) {
	addr, err := askUint16("Starting address")
	if err != nil {
		return nil, fmt.Errorf("ask address: %w", err)
	}

	quantity, err := askUint16("Quantity")
	if err != nil {
		return nil, fmt.Errorf("ask quantity: %w", err)
	}

	inputs, err := a.client.ReadDiscreteInputs(addr, quantity)
	if err != nil {
		return nil, fmt.Errorf(
			"could not read %d discrete inputs at %d: %w",
			quantity, addr, err,
		)
	}

	return inputs, nil
}

// 0x03 Read holding registers
func (a *Application) ReadHoldingRegisters() ([]uint16, error) {
	addr, err := askUint16("Starting address")
	if err != nil {
		return nil, fmt.Errorf("ask address: %w", err)
	}

	quantity, err := askUint16("Quantity")
	if err != nil {
		return nil, fmt.Errorf("ask quantity: %w", err)
	}

	registers, err := a.client.ReadRegisters(addr, quantity, modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, fmt.Errorf("could not read %d hodling registers at %d: %w",
			quantity, addr, err,
		)
	}

	return registers, nil
}

// 0x04 Read input registers
func (a *Application) ReadInputRegisters() ([]uint16, error) {
	addr, err := askUint16("Starting address")
	if err != nil {
		return nil, fmt.Errorf("ask address: %w", err)
	}

	quantity, err := askUint16("Quantity")
	if err != nil {
		return nil, fmt.Errorf("ask quantity: %w", err)
	}

	registers, err := a.client.ReadRegisters(addr, quantity, modbus.INPUT_REGISTER)
	if err != nil {
		return nil, fmt.Errorf("could not read %d input registers at %d: %w",
			quantity, addr, err,
		)
	}

	return registers, nil
}

// 0x05 Write single coil
func (a *Application) WriteSingleCoil() error {
	addr, err := askUint16("Address")
	if err != nil {
		return fmt.Errorf("ask address: %w", err)
	}

	coil, err := askBool("Coil")
	if err != nil {
		return fmt.Errorf("ask coil: %w", err)
	}

	if err := a.client.WriteCoil(addr, coil); err != nil {
		return fmt.Errorf("could not write %v coil at %d: %w", coil, addr, err)
	}

	return nil
}

// 0x06 Read single coil
func (a *Application) WriteSingleRegister() error {
	addr, err := askUint16("Starting address")
	if err != nil {
		return fmt.Errorf("ask address: %w", err)
	}

	value, err := askUint16("Value")
	if err != nil {
		return fmt.Errorf("ask value: %w", err)
	}

	if err := a.client.WriteRegister(addr, value); err != nil {
		return fmt.Errorf("could not write value %d at %d: %w", addr, value, err)
	}

	return nil
}

// 0x10 Write mutliple registers
func (a *Application) WriteMultipleRegisters() error {
	addr, err := askUint16("Starting address")
	if err != nil {
		return fmt.Errorf("ask address: %w", err)
	}

	quantity, err := askUint16("Quantity")
	if err != nil {
		return fmt.Errorf("ask quantity: %w", err)
	}

	values := make([]uint16, 0, quantity)
	for i := 0; i < int(quantity); i++ {
		value, err := askUint16(fmt.Sprintf("Value[%d]", i))
		if err != nil {
			return fmt.Errorf("ask value[%d]: %w", i, err)
		}

		values = append(values, value)
	}

	if err := a.client.WriteRegisters(addr, values); err != nil {
		return fmt.Errorf("could not write %v values %v at %v: %w",
			quantity, values, addr, err)
	}

	return nil
}

// 0x0F Write multiple coils
func (a *Application) WriteMultipleCoils() error {
	addr, err := askUint16("Starting address")
	if err != nil {
		return fmt.Errorf("ask address: %w", err)
	}

	quantity, err := askUint16("Quantity")
	if err != nil {
		return fmt.Errorf("ask quantity: %w", err)
	}

	coils := make([]bool, 0, quantity)
	for i := 0; i < int(quantity); i++ {
		coil, err := askBool(fmt.Sprintf("Coil[%d]", i))
		if err != nil {
			return fmt.Errorf("ask coil[%d]: %w", i, err)
		}

		coils = append(coils, coil)
	}

	if err := a.client.WriteCoils(addr, coils); err != nil {
		return fmt.Errorf("could not write %v coils %v at %v: %w",
			quantity, coils, addr, err)
	}

	return nil
}

func askUint16(message string) (uint16, error) {
	fmt.Printf("%s: ", message)

	var input string
	fmt.Scanln(&input)

	inputU64, err := strconv.ParseUint(input, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("error parsing uint16 %s: %w", input, err)
	}

	return uint16(inputU64), nil
}

func askBool(message string) (bool, error) {
	fmt.Printf("%s: ", message)

	var input string
	fmt.Scanln(&input)

	if input == "true" || input == "1" || input == "t" || input == "T" {
		return true, nil
	}

	if input == "false" || input == "0" || input == "f" || input == "F" {
		return false, nil
	}

	return false, errors.New("unsupported boolean format")
}

func askString(message string) string {
	fmt.Printf("%s: ", message)
	var input string
	fmt.Scanln(&input)
	return input
}

func main() {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL: "tcp://localhost:5502",
	})

	if err != nil {
		panic(fmt.Errorf("failed to create modbus client: %w", err))
	}

	if err := client.Open(); err != nil {
		panic(fmt.Errorf("could not connect: %w", err))
	}

	app := &Application{
		client: client,
	}

	handlers := map[string]func(){
		"1": func() {
			coils, err := app.ReadCoils()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(coils)
		},

		"2": func() {
			values, err := app.ReadDiscreteInputs()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(values)
		},

		"3": func() {
			registers, err := app.ReadHoldingRegisters()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(registers)
		},

		"4": func() {
			registers, err := app.ReadInputRegisters()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(registers)
		},

		"5": func() {
			if err := app.WriteSingleCoil(); err != nil {
				fmt.Println(err)
			}
		},

		"6": func() {
			if err := app.WriteSingleRegister(); err != nil {
				fmt.Println(err)
			}
		},

		"7": func() {
			if err := app.WriteMultipleRegisters(); err != nil {
				fmt.Println(err)
			}
		},

		"8": func() {
			if err := app.WriteMultipleCoils(); err != nil {
				fmt.Println(err)
			}
		},
	}

	for {
		fmt.Println("1: 0x01 Read coils")
		fmt.Println("2: 0x02 Read discrete inputs")
		fmt.Println("3: 0x03 Read holding registers")
		fmt.Println("4: 0x04 Read input registers")
		fmt.Println("5: 0x05 Write single coil")
		fmt.Println("6: 0x06 Write single register")
		fmt.Println("7: 0x10 Write multiple registers")
		fmt.Println("8: 0x0F Write multiple coils")
		fmt.Println()
		fmt.Println("0: Exit")

		function := askString("Function")
		handler, ok := handlers[function]
		if !ok {
			if function == "0" {
				break
			}

			continue
		}

		handler()
	}

	if err := client.Close(); err != nil {
		panic(fmt.Errorf("close client: %w", err))
	}
}
