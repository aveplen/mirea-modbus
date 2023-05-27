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

type DialogModel interface {
	ReadCoils(addr uint16, cnt int) ([]bool, error)
	ReadDiscreteInputs(addr uint16, cnt int) ([]bool, error)
	ReadHoldingRegisters(addr uint16, cnt int) ([]uint16, error)
	ReadInputRegisters(addr uint16, cnt int) ([]uint16, error)
	WriteSingleCoil(addr uint16, value bool) error
	WriteSingleRegister(addr uint16, value uint16) error
	WriteMultipleRegisters(addr uint16, values []uint16) error
	WriteMultipleCoils(addr uint16, values []bool) error
}

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

type DialogController struct {
	model DialogModel

	dialog           *walk.Dialog
	errEdit          *walk.TextEdit
	addrEdit         *walk.TextEdit
	hexAddrCheckBox  *walk.CheckBox
	hexInputCheckBox *walk.CheckBox
	binaryInputEdit  *walk.TextEdit
	inputEdit        *walk.TextEdit
	cntEdit          *walk.TextEdit
	resultEdit       *walk.TextEdit
}

func (c *DialogController) Close() {
	c.dialog.Close(0)
}

func (c *DialogController) ReadCoils() {
	addr, cnt, ok := c.addrCnt()
	if !ok {
		return
	}

	coils, err := c.model.ReadCoils(addr, cnt)
	if err != nil {
		c.setError(err)
		c.resultFail()
		return
	}

	c.resultBools(coils)
	c.clearError()
}

func (c *DialogController) ReadDiscreteInputs() {
	addr, cnt, ok := c.addrCnt()
	if !ok {
		return
	}

	inputs, err := c.model.ReadDiscreteInputs(addr, int(cnt))
	if err != nil {
		c.setError(err)
		c.resultFail()
		return
	}

	c.resultBools(inputs)
	c.clearError()
}

func (c *DialogController) ReadHoldingRegisters() {
	addr, cnt, ok := c.addrCnt()
	if !ok {
		return
	}

	registers, err := c.model.ReadHoldingRegisters(addr, int(cnt))
	if err != nil {
		c.setError(err)
		c.resultFail()
		return
	}

	c.resultUints(registers)
	c.clearError()
}

func (c *DialogController) ReadInputRegisters() {
	addr, cnt, ok := c.addrCnt()
	if !ok {
		return
	}

	registers, err := c.model.ReadInputRegisters(addr, int(cnt))
	if err != nil {
		c.setError(err)
		c.resultFail()
		return
	}

	c.resultUints(registers)
	c.clearError()
}

func (c *DialogController) WriteSingleCoil() {
	addr, ok := c.addr()
	if !ok {
		return
	}

	input, ok := c.inputBool()
	if !ok {
		return
	}

	if err := c.model.WriteSingleCoil(addr, input); err != nil {
		c.setError(err)
		c.resultFail()
		return
	}

	c.resultSuccess()
	c.clearError()
}

func (c *DialogController) WriteSingleRegister() {
	addr, ok := c.addr()
	if !ok {
		return
	}

	input, ok := c.inputUint()
	if !ok {
		return
	}

	if err := c.model.WriteSingleRegister(addr, input); err != nil {
		c.setError(err)
		c.resultFail()
		return
	}

	c.resultSuccess()
	c.clearError()
}

func (c *DialogController) WriteMultipleRegisters() {
	addr, ok := c.addr()
	if !ok {
		return
	}

	inputs, ok := c.inputUints()
	if !ok {
		return
	}

	if err := c.model.WriteMultipleRegisters(addr, inputs); err != nil {
		c.setError(err)
		c.resultFail()
		return
	}

	c.resultSuccess()
	c.clearError()
}

func (c *DialogController) WriteMultipleCoils() {
	addr, ok := c.addr()
	if !ok {
		return
	}

	inputs, ok := c.inputBools()
	if !ok {
		return
	}

	if err := c.model.WriteMultipleCoils(addr, inputs); err != nil {
		c.setError(err)
		c.resultFail()
		return
	}

	c.resultSuccess()
	c.clearError()
}

func (c *DialogController) addr() (uint16, bool) {
	addrParser := parseUint16
	if c.hexAddrCheckBox.Checked() {
		addrParser = parseHex
	}

	addr, err := addrParser(c.addrEdit.Text())
	if err != nil {
		c.setError(err)
		return 0, false
	}

	c.clearError()
	return addr, true
}

func (c *DialogController) cnt() (int, bool) {
	cnt, err := parseInt(c.cntEdit.Text())
	if err != nil {
		c.setError(err)
		return 0, false
	}

	c.clearError()
	return cnt, true
}

func (c *DialogController) inputBool() (bool, bool) {
	parsed, err := parseBool(c.inputEdit.Text())
	if err != nil {
		c.setError(err)
		return false, false
	}

	c.clearError()
	return parsed, true
}

func (c *DialogController) inputBools() ([]bool, bool) {
	cnt, ok := c.cnt()
	if !ok {
		return nil, false
	}

	input := make([]bool, 0, cnt)
	for _, chunk := range strings.Split(c.inputEdit.Text(), ", ") {
		parsed, err := parseBool(chunk)

		if err != nil {
			c.setError(err)
			return nil, false
		}

		input = append(input, parsed)
	}

	c.clearError()
	return input, true
}

func (c *DialogController) inputUint() (uint16, bool) {
	inputParser := parseUint16
	if c.hexInputCheckBox.Checked() {
		inputParser = parseHex
	}

	parsed, err := inputParser(c.inputEdit.Text())
	if err != nil {
		c.setError(err)
		return 0, false
	}

	c.clearError()
	return parsed, true
}

func (c *DialogController) inputUints() ([]uint16, bool) {
	cnt, ok := c.cnt()
	if !ok {
		return nil, false
	}

	inputParser := parseUint16
	if c.hexInputCheckBox.Checked() {
		inputParser = parseHex
	}

	input := make([]uint16, 0, cnt)
	for _, chunk := range strings.Split(c.inputEdit.Text(), ", ") {
		parsed, err := inputParser(chunk)

		if err != nil {
			c.setError(err)
			return nil, false
		}

		input = append(input, parsed)
	}

	c.clearError()
	return input, true
}

func (c *DialogController) addrCnt() (uint16, int, bool) {
	addr, ok := c.addr()
	if !ok {
		return 0, 0, false
	}

	cnt, ok := c.cnt()
	if !ok {
		return 0, 0, false
	}

	c.clearError()
	return addr, cnt, true
}

func (c *DialogController) resultSuccess() {
	c.resultEdit.SetText("Success")
}

func (c *DialogController) resultFail() {
	c.resultEdit.SetText("Fail")
}

func (c *DialogController) resultUints(values []uint16) {
	c.resultEdit.SetText(fmt.Sprintf("%v", values))
}

func (c *DialogController) resultBools(values []bool) {
	c.resultEdit.SetText(fmt.Sprintf("%v", values))
}

func (c *DialogController) setError(err error) {
	c.errEdit.SetText(err.Error())
}

func (c *DialogController) clearError() {
	c.errEdit.SetText("")
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

func DialogView(window *walk.MainWindow, model DialogModel, dialogType DialogType) func() {
	controller := &DialogController{
		model: model,
	}

	mainButtonFunction := map[DialogType]func(){
		DialogTypeReadCoils:              controller.ReadCoils,
		DialogTypeReadDiscreteInputs:     controller.ReadDiscreteInputs,
		DialogTypeReadHoldingRegisters:   controller.ReadHoldingRegisters,
		DialogTypeReadInputRegisters:     controller.ReadInputRegisters,
		DialogTypeWriteSingleCoil:        controller.WriteSingleCoil,
		DialogTypeWriteSingleRegister:    controller.WriteSingleRegister,
		DialogTypeWriteMultipleRegisters: controller.WriteMultipleRegisters,
		DialogTypeWriteMultipleCoils:     controller.WriteMultipleCoils,
	}

	return func() {
		d.Dialog{
			AssignTo: &controller.dialog,
			Title:    dialogTitles[dialogType],
			MinSize:  d.Size{Width: 250, Height: 200},
			Layout:   d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
			Children: func() []d.Widget {

				widgets := make([]d.Widget, 0)

				widgets = append(widgets, d.GroupBox{
					Title:  "Starting address (hex or decimal)",
					Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
					Children: []d.Widget{
						d.TextEdit{AssignTo: &controller.addrEdit, Text: "0x01"},
						d.CheckBox{AssignTo: &controller.hexAddrCheckBox, Checked: true, Text: "Hexadecimal number"},
					},
				})

				if renderAmount[dialogType] {
					widgets = append(widgets, d.GroupBox{
						Title:  "Amount (decimal)",
						Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
						Children: []d.Widget{
							d.TextEdit{AssignTo: &controller.cntEdit, Text: inputDefaultText[dialogType]},
						},
					})
				}

				if renderBinaryInput[dialogType] {
					widgets = append(widgets, d.GroupBox{
						Title:  "Input (example: 'true, false')",
						Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
						Children: []d.Widget{
							d.TextEdit{AssignTo: &controller.binaryInputEdit, Text: inputDefaultText[dialogType]},
						},
					})
				}

				if renderInput[dialogType] {
					widgets = append(widgets, d.GroupBox{
						Title:  "Input (example: '213, 45' or '0x15, 0x1')",
						Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
						Children: []d.Widget{
							d.TextEdit{AssignTo: &controller.inputEdit, Text: inputDefaultText[dialogType]},
							d.CheckBox{AssignTo: &controller.hexInputCheckBox, Checked: true, Text: "Hexadecimal number"},
						},
					})
				}

				widgets = append(widgets, d.GroupBox{
					Title:  "Result",
					Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
					Children: []d.Widget{
						d.TextEdit{AssignTo: &controller.resultEdit, Enabled: false},
					},
				})

				widgets = append(widgets, d.HSplitter{
					Children: []d.Widget{
						d.PushButton{Text: mainButtonCaption[dialogType], OnClicked: mainButtonFunction[dialogType]},
						d.PushButton{Text: "Cancel", OnClicked: controller.Close},
					},
				})

				widgets = append(widgets, d.TextEdit{
					AssignTo:  &controller.errEdit,
					TextColor: walk.RGB(255, 0, 0),
					ReadOnly:  true,
					Text:      "",
				})

				return widgets
			}(),
		}.Run(window)
	}
}
