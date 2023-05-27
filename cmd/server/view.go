package main

import (
	"fmt"

	"github.com/lxn/walk"
	d "github.com/lxn/walk/declarative"
)

type CoilView struct {
	Index   int
	Address uint16
	Value   bool
}

type CoilsModel struct {
	walk.TableModelBase
	items []*CoilView
}

func NewCoilsModel(coils []Coil) *CoilsModel {
	m := new(CoilsModel)
	m.items = make([]*CoilView, len(coils))
	for i := range m.items {
		m.items[i] = &CoilView{
			Index:   i,
			Address: coils[i].addr,
			Value:   coils[i].value,
		}
	}
	m.PublishRowsReset()
	return m
}

func (m *CoilsModel) RowCount() int {
	return len(m.items)
}

func (m *CoilsModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1, 2:
		return item.Address

	case 3:
		return item.Value
	}

	panic("unexpected col")
}

type RegisterView struct {
	Index   int
	Address uint16
	Value   uint16
}

type RegistersModel struct {
	walk.TableModelBase
	items []*RegisterView
}

func NewRegistersModel(registers []Register) *RegistersModel {
	m := new(RegistersModel)
	m.ResetRows(registers)
	return m
}

func (m *RegistersModel) ResetRows(registers []Register) {
	m.items = make([]*RegisterView, len(registers))
	for i := range m.items {
		m.items[i] = &RegisterView{
			Index:   i,
			Address: registers[i].addr,
			Value:   registers[i].value,
		}
	}
	m.PublishRowsReset()
}

func (m *RegistersModel) RowCount() int {
	return len(m.items)
}

func (m *RegistersModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1, 2:
		return item.Address

	case 3, 4:
		return item.Value
	}

	panic("unexpected col coils")
}

type MainModel interface {
	StartServer() bool
	StopServer() bool
}

type ViewController struct {
	MainWindow *d.MainWindow
	model      MainModel

	discreteInputsModel   *CoilsModel
	coilsModel            *CoilsModel
	inputRegistersModel   *RegistersModel
	holdingRegistersModel *RegistersModel

	discreteInputsView   *walk.TableView
	coilsView            *walk.TableView
	inputRegisterView    *walk.TableView
	holdingRegistersView *walk.TableView
	startServerButton    *walk.PushButton
	stopServerButton     *walk.PushButton

	AppendLog func(value string)
}

func (v *ViewController) UpdateDiscreteInputs(change CoilChange) {
	for _, item := range v.discreteInputsModel.items {
		if item.Address == change.addr {
			item.Value = change.to
			break
		}
	}

	v.coilsModel.PublishRowsReset()
	v.coilsView.Invalidate()
}

func (v *ViewController) UpdateCoils(change CoilChange) {
	for _, item := range v.coilsModel.items {
		if item.Address == change.addr {
			item.Value = change.to
			break
		}
	}

	v.coilsModel.PublishRowsReset()
	v.coilsView.Invalidate()
}

func (v *ViewController) UpdateInputRegisters(change RegisterChange) {
	for _, item := range v.inputRegistersModel.items {
		if item.Address == change.addr {
			item.Value = change.to
			break
		}
	}

	v.inputRegistersModel.PublishRowsReset()
	v.inputRegisterView.Invalidate()
}

func (v *ViewController) UpdateHoldingRegisters(change RegisterChange) {
	for _, item := range v.holdingRegistersModel.items {
		if item.Address == change.addr {
			item.Value = change.to
			break
		}
	}

	v.holdingRegistersModel.PublishRowsReset()
	v.holdingRegistersView.Invalidate()
}

func (v *ViewController) StartServer() {
	if v.model.StartServer() {
		v.startServerButton.Button.SetEnabled(false)
		v.stopServerButton.Button.SetEnabled(true)
	}
}

func (v *ViewController) StopServer() {
	if v.model.StopServer() {
		v.startServerButton.Button.SetEnabled(true)
		v.stopServerButton.Button.SetEnabled(false)
	}
}

func NewView(seed Dump, model MainModel) ViewController {
	view := ViewController{
		discreteInputsModel:   NewCoilsModel(seed.DiscreteInputs),
		coilsModel:            NewCoilsModel(seed.Coils),
		inputRegistersModel:   NewRegistersModel(seed.InputRegisters),
		holdingRegistersModel: NewRegistersModel(seed.HoldingRegisters),
	}

	lv := NewLogView()
	view.MainWindow = &d.MainWindow{
		Title:  "Modbus server (slave)",
		Size:   d.Size{Width: 550, Height: 900},
		Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
		Children: []d.Widget{
			d.HSplitter{
				Children: []d.Widget{
					d.PushButton{
						AssignTo:  &view.startServerButton,
						Text:      "Start server",
						OnClicked: view.StartServer,
					},

					d.PushButton{
						AssignTo:  &view.stopServerButton,
						Text:      "Stop server",
						OnClicked: view.StopServer,
						Enabled:   false,
					},
				},
			},

			d.Label{Text: "Coils:"},
			d.TableView{
				AssignTo:         &view.coilsView,
				Model:            view.coilsModel,
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				MaxSize:          d.Size{Height: 140},
				Columns: []d.TableViewColumn{
					{Title: "#", FormatFunc: numberDecimalFormat},
					{Title: "Address (dec)", FormatFunc: numberDecimalFormat},
					{Title: "Address (hex)", FormatFunc: numberHexFormat},
					{Title: "Coil", FormatFunc: boolFormat},
				},
			},

			d.Label{Text: "Discrete inputs:"},
			d.TableView{
				AssignTo:         &view.discreteInputsView,
				Model:            view.discreteInputsModel,
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				MaxSize:          d.Size{Height: 140},
				Columns: []d.TableViewColumn{
					{Title: "#", FormatFunc: numberDecimalFormat},
					{Title: "Address (dec)", FormatFunc: numberDecimalFormat},
					{Title: "Address (hex)", FormatFunc: numberHexFormat},
					{Title: "Coil", FormatFunc: boolFormat},
				},
			},

			d.Label{Text: "Input registers:"},
			d.TableView{
				AssignTo:         &view.inputRegisterView,
				Model:            view.inputRegistersModel,
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				MaxSize:          d.Size{Height: 180},
				Columns: []d.TableViewColumn{
					{Title: "#", FormatFunc: numberDecimalFormat},
					{Title: "Address (dec)", FormatFunc: numberDecimalFormat},
					{Title: "Address (hex)", FormatFunc: numberHexFormat},
					{Title: "Value (dec)", FormatFunc: numberDecimalFormat},
					{Title: "Value (hex)", FormatFunc: numberHexFormat},
				},
			},

			d.Label{Text: "Holding registers:"},
			d.TableView{
				AssignTo:         &view.holdingRegistersView,
				Model:            view.holdingRegistersModel,
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				MaxSize:          d.Size{Height: 140},
				Columns: []d.TableViewColumn{
					{Title: "#", FormatFunc: numberDecimalFormat},
					{Title: "Address (dec)", FormatFunc: numberDecimalFormat},
					{Title: "Address (hex)", FormatFunc: numberHexFormat},
					{Title: "Value (dec)", FormatFunc: numberDecimalFormat},
					{Title: "Value (hex)", FormatFunc: numberHexFormat},
				},
			},

			lv.TextEdit,
		},
	}

	return view
}

func numberDecimalFormat(value interface{}) string {
	return fmt.Sprintf("%d", value)
}

func numberHexFormat(value interface{}) string {
	return fmt.Sprintf("0x%X", value)
}

func boolFormat(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
