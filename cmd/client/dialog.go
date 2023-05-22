package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lxn/walk"
	d "github.com/lxn/walk/declarative"
)

type DialogType int

const (
	DialogTypeReadCoils DialogType = iota + 1
	DialogTypeReadDiscreteInputs
	DialogTypeReadHoldingRegisters
	DialogTypeReadInputRegisters
	DialogTypeWriteSingleCoil
	DialogTypeWriteSingleRegister
	DialogTypeWriteMultipleRegisters
	DialogTypeWriteMultipleCoils
)

var (
	dialogTitles = map[DialogType]string{
		DialogTypeReadCoils:              "Read coils 0x01",
		DialogTypeReadDiscreteInputs:     "Read discrete inputs 0x02",
		DialogTypeReadHoldingRegisters:   "Read holding registers 0x03",
		DialogTypeReadInputRegisters:     "Read input registesrs 0x04",
		DialogTypeWriteSingleCoil:        "Write single coild 0x05",
		DialogTypeWriteSingleRegister:    "Write single register 0x06",
		DialogTypeWriteMultipleRegisters: "Write multiple registers 0x10",
		DialogTypeWriteMultipleCoils:     "Write multiple coils 0x0F",
	}

	renderAmount = map[DialogType]bool{
		DialogTypeReadCoils:              true,
		DialogTypeReadDiscreteInputs:     true,
		DialogTypeReadHoldingRegisters:   true,
		DialogTypeReadInputRegisters:     true,
		DialogTypeWriteSingleCoil:        false,
		DialogTypeWriteSingleRegister:    false,
		DialogTypeWriteMultipleRegisters: true,
		DialogTypeWriteMultipleCoils:     true,
	}

	renderBinaryInput = map[DialogType]bool{
		DialogTypeReadCoils:              false,
		DialogTypeReadDiscreteInputs:     false,
		DialogTypeReadHoldingRegisters:   false,
		DialogTypeReadInputRegisters:     false,
		DialogTypeWriteSingleCoil:        true,
		DialogTypeWriteSingleRegister:    false,
		DialogTypeWriteMultipleRegisters: false,
		DialogTypeWriteMultipleCoils:     true,
	}

	renderInput = map[DialogType]bool{
		DialogTypeReadCoils:              false,
		DialogTypeReadDiscreteInputs:     false,
		DialogTypeReadHoldingRegisters:   false,
		DialogTypeReadInputRegisters:     false,
		DialogTypeWriteSingleCoil:        false,
		DialogTypeWriteSingleRegister:    true,
		DialogTypeWriteMultipleRegisters: true,
		DialogTypeWriteMultipleCoils:     false,
	}

	inputDefaultText = map[DialogType]string{
		DialogTypeReadCoils:              "",
		DialogTypeReadDiscreteInputs:     "",
		DialogTypeReadHoldingRegisters:   "",
		DialogTypeReadInputRegisters:     "",
		DialogTypeWriteSingleCoil:        "true",
		DialogTypeWriteSingleRegister:    "0x123",
		DialogTypeWriteMultipleRegisters: "0x123, 0x456",
		DialogTypeWriteMultipleCoils:     "true, false",
	}

	mainButtonCaption = map[DialogType]string{
		DialogTypeReadCoils:              "Read",
		DialogTypeReadDiscreteInputs:     "Read",
		DialogTypeReadHoldingRegisters:   "Read",
		DialogTypeReadInputRegisters:     "Read",
		DialogTypeWriteSingleCoil:        "Write",
		DialogTypeWriteSingleRegister:    "Write",
		DialogTypeWriteMultipleRegisters: "Write",
		DialogTypeWriteMultipleCoils:     "Write",
	}
)

type DialogCallbacks struct {
	readCoils              func(addr uint16, cnt int) ([]bool, error)
	readDiscreteInputs     func(addr uint16, cnt int) ([]bool, error)
	readHoldingRegisters   func(addr uint16, cnt int) ([]uint16, error)
	readInputRegisters     func(addr uint16, cnt int) ([]uint16, error)
	writeSingleCoil        func(addr uint16, value bool) error
	writeSingleRegister    func(addr uint16, value uint16) error
	writeMultipleRegisters func(addr uint16, values []uint16) error
	writeMultipleCoils     func(addr uint16, values []bool) error
}

type Dialog struct {
	callbacks  DialogCallbacks
	dialogType DialogType

	dialog   *walk.Dialog
	errEdit  *walk.TextEdit
	addrEdit *walk.TextEdit

	hexAddrCheckBox  *walk.CheckBox
	hexInputCheckBox *walk.CheckBox

	binaryInputEdit *walk.TextEdit
	inputEdit       *walk.TextEdit

	cntEdit    *walk.TextEdit
	resultEdit *walk.TextEdit

	mainButtonFunction map[DialogType]func()
}

func NewDialog(
	callbacks DialogCallbacks,
	dialogType DialogType,
) *Dialog {

	d := &Dialog{
		callbacks:  callbacks,
		dialogType: dialogType,
	}

	d.mainButtonFunction = map[DialogType]func(){
		DialogTypeReadCoils:              d.readCoils,
		DialogTypeReadDiscreteInputs:     d.readDiscreteInputs,
		DialogTypeReadHoldingRegisters:   d.readHoldingRegisters,
		DialogTypeReadInputRegisters:     d.readInputRegisters,
		DialogTypeWriteSingleCoil:        d.writeSingleCoil,
		DialogTypeWriteSingleRegister:    d.writeSingleRegister,
		DialogTypeWriteMultipleRegisters: d.writeMultipleRegisters,
		DialogTypeWriteMultipleCoils:     d.writeMultipleCoils,
	}

	return d
}

func (g *Dialog) Close() {
	g.dialog.Close(0)
}

func (g *Dialog) readCoils() {
	addrParser := parseUint16
	if g.hexAddrCheckBox.Checked() {
		addrParser = parseHex
	}

	addr, err := addrParser(g.addrEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	cnt, err := parseInt(g.cntEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	coils, err := g.callbacks.readCoils(addr, int(cnt))
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	g.resultEdit.SetText(writeBools(coils))
}

func (g *Dialog) readDiscreteInputs() {
	addrParser := parseUint16
	if g.hexAddrCheckBox.Checked() {
		addrParser = parseHex
	}

	addr, err := addrParser(g.addrEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	cnt, err := parseInt(g.cntEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	inputs, err := g.callbacks.readDiscreteInputs(addr, int(cnt))
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	g.resultEdit.SetText(writeBools(inputs))
}

func (g *Dialog) readHoldingRegisters() {
	addrParser := parseUint16
	if g.hexAddrCheckBox.Checked() {
		addrParser = parseHex
	}

	addr, err := addrParser(g.addrEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	cnt, err := parseInt(g.cntEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	registers, err := g.callbacks.readHoldingRegisters(addr, int(cnt))
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	g.resultEdit.SetText(writeUints(registers))
}

func (g *Dialog) readInputRegisters() {
	addrParser := parseUint16
	if g.hexAddrCheckBox.Checked() {
		addrParser = parseHex
	}

	addr, err := addrParser(g.addrEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	cnt, err := parseInt(g.cntEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	registers, err := g.callbacks.readInputRegisters(addr, int(cnt))
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	g.resultEdit.SetText(writeUints(registers))
}

func (g *Dialog) writeSingleCoil() {
	addrParser := parseUint16
	if g.hexAddrCheckBox.Checked() {
		addrParser = parseHex
	}

	addr, err := addrParser(g.addrEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	input, err := parseBool(g.binaryInputEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	if err := g.callbacks.writeSingleCoil(addr, input); err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	g.resultEdit.SetText("Success")
}

func (g *Dialog) writeSingleRegister() {
	addrParser := parseUint16
	if g.hexAddrCheckBox.Checked() {
		addrParser = parseHex
	}

	addr, err := addrParser(g.addrEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	inputParser := parseUint16
	if g.hexInputCheckBox.Checked() {
		inputParser = parseHex
	}

	input, err := inputParser(g.inputEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	if err := g.callbacks.writeSingleRegister(addr, input); err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	g.resultEdit.SetText("Success")
}

func (g *Dialog) writeMultipleRegisters() {
	addrParser := parseUint16
	if g.hexAddrCheckBox.Checked() {
		addrParser = parseHex
	}

	addr, err := addrParser(g.addrEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	cnt, err := parseInt(g.cntEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	inputParser := parseUint16
	if g.hexInputCheckBox.Checked() {
		inputParser = parseHex
	}

	input := make([]uint16, 0, cnt)
	for _, chunk := range strings.Split(g.inputEdit.Text(), ", ") {
		parsed, err := inputParser(chunk)

		if err != nil {
			g.errEdit.SetText(err.Error())
			return
		}

		input = append(input, parsed)
	}

	if err := g.callbacks.writeMultipleRegisters(addr, input); err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	g.resultEdit.SetText("Success")
}

func (g *Dialog) writeMultipleCoils() {
	addrParser := parseUint16
	if g.hexAddrCheckBox.Checked() {
		addrParser = parseHex
	}

	addr, err := addrParser(g.addrEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	cnt, err := parseInt(g.cntEdit.Text())
	if err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	input := make([]bool, 0, cnt)
	for _, chunk := range strings.Split(g.inputEdit.Text(), ", ") {
		parsed, err := parseBool(chunk)

		if err != nil {
			g.errEdit.SetText(err.Error())
			return
		}

		input = append(input, parsed)
	}

	if err := g.callbacks.writeMultipleCoils(addr, input); err != nil {
		g.errEdit.SetText(err.Error())
		return
	}

	g.resultEdit.SetText("Success")
}

func (g *Dialog) Render() d.Dialog {
	return d.Dialog{
		AssignTo: &g.dialog,
		Title:    dialogTitles[g.dialogType],
		MinSize:  d.Size{Width: 250, Height: 200},
		Layout:   d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
		Children: func() []d.Widget {

			widgets := make([]d.Widget, 0)

			widgets = append(widgets, d.GroupBox{
				Title:  "Starting address (hex or decimal)",
				Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
				Children: []d.Widget{
					d.TextEdit{AssignTo: &g.addrEdit, Text: "0x01"},
					d.CheckBox{AssignTo: &g.hexAddrCheckBox, Checked: true, Text: "Hexadecimal number"},
				},
			})

			if renderAmount[g.dialogType] {
				widgets = append(widgets, d.GroupBox{
					Title:  "Amount (decimal)",
					Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
					Children: []d.Widget{
						d.TextEdit{AssignTo: &g.cntEdit, Text: inputDefaultText[g.dialogType]},
					},
				})
			}

			if renderBinaryInput[g.dialogType] {
				widgets = append(widgets, d.GroupBox{
					Title:  "Input (example: 'true, false')",
					Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
					Children: []d.Widget{
						d.TextEdit{AssignTo: &g.binaryInputEdit, Text: inputDefaultText[g.dialogType]},
					},
				})
			}

			if renderInput[g.dialogType] {
				widgets = append(widgets, d.GroupBox{
					Title:  "Input (example: '213, 45' or '0x15, 0x1')",
					Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
					Children: []d.Widget{
						d.TextEdit{AssignTo: &g.inputEdit, Text: inputDefaultText[g.dialogType]},
						d.CheckBox{AssignTo: &g.hexInputCheckBox, Checked: true, Text: "Hexadecimal number"},
					},
				})
			}

			widgets = append(widgets, d.GroupBox{
				Title:  "Result",
				Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
				Children: []d.Widget{
					d.TextEdit{AssignTo: &g.resultEdit, Enabled: false},
				},
			})

			widgets = append(widgets, d.HSplitter{
				Children: []d.Widget{
					d.PushButton{Text: mainButtonCaption[g.dialogType], OnClicked: g.mainButtonFunction[g.dialogType]},
					d.PushButton{Text: "Cancel", OnClicked: g.Close},
				},
			})

			widgets = append(widgets, d.TextEdit{
				AssignTo:  &g.errEdit,
				TextColor: walk.RGB(255, 0, 0),
				ReadOnly:  true,
				Text:      "",
			})

			return widgets
		}(),
	}
}

func parseHex(input string) (uint16, error) {
	if input[:2] != "0x" {
		return 0, errors.New("hex should start with '0x'")
	}

	u64, err := strconv.ParseUint(input[2:], 16, 64)
	if err != nil {
		return 0, errors.New("could not parse hex input")
	}

	if u64 > 0xFFFF {
		return 0, errors.New("value must be less or equal then 0xFFFF")
	}

	return uint16(u64), nil
}

func parseUint16(input string) (uint16, error) {
	u64, err := strconv.ParseUint(input, 16, 64)
	if err != nil {
		return 0, errors.New("could not parse decimal input")
	}

	if u64 > 0xFFFF {
		return 0, fmt.Errorf("value must be less or equal then %d", 0xFFFF)
	}

	return uint16(u64), nil
}

func parseInt(input string) (int, error) {
	i, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("could not parse integer")
	}
	return i, nil
}

func parseBool(input string) (bool, error) {
	if input == "true" {
		return true, nil
	}

	if input == "false" {
		return false, nil
	}

	return false, fmt.Errorf("could not parse bool")
}

func writeBools(bools []bool) string {
	return fmt.Sprintf("%v", bools)
}

func writeUints(uints []uint16) string {
	return fmt.Sprintf("%v", uints)
}
