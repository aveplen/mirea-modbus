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

	case 1:
		return item.Address

	case 2:
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

	case 1:
		return item.Address

	case 2, 3:
		return item.Value
	}

	panic("unexpected col")
}

type View struct {
	CoilsUpdateCallback          func(change CoilChange)
	RegistersRenderCallback      func(change RegisterChange)
	InputRegistersRenderCallback func(change RegisterChange)
	MainWindow                   *d.MainWindow
}

func NewView(
	seed Dump,
	startServerCallback func() bool,
	stopServerCallback func() bool,
) View {
	coilsModel := NewCoilsModel(seed.Coils)
	var coilsView *walk.TableView

	registersModel := NewRegistersModel(seed.Registers)
	var registersView *walk.TableView

	inputRegistersModel := NewRegistersModel(seed.InputRegisters)
	var inputRegisterView *walk.TableView

	var startServerButton *walk.PushButton
	var stopServerButton *walk.PushButton

	return View{
		CoilsUpdateCallback: func(change CoilChange) {
			for _, item := range coilsModel.items {
				if item.Address == change.addr {
					item.Value = change.to
					break
				}
			}
			coilsModel.PublishRowsReset()
			coilsView.SetEnabled(false)
			coilsView.SetEnabled(true)
		},

		RegistersRenderCallback: func(change RegisterChange) {
			for _, item := range registersModel.items {
				if item.Address == change.addr {
					item.Value = change.to
					break
				}
			}
			registersModel.PublishRowsReset()
			registersView.SetEnabled(false)
			registersView.SetEnabled(true)
		},

		InputRegistersRenderCallback: func(change RegisterChange) {
			for _, item := range inputRegistersModel.items {
				if item.Address == change.addr {
					item.Value = change.to
					break
				}
			}
			inputRegistersModel.PublishRowsReset()
			inputRegisterView.SetEnabled(false)
			inputRegisterView.SetEnabled(true)
		},

		MainWindow: &d.MainWindow{
			Title:  "Modbus server (slave)",
			Size:   d.Size{Width: 500, Height: 800},
			Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
			Children: []d.Widget{
				d.HSplitter{
					Children: []d.Widget{
						d.PushButton{
							AssignTo: &startServerButton,
							Text:     "Start server",
							OnClicked: func() {
								if startServerCallback() {
									startServerButton.Button.SetEnabled(false)
									stopServerButton.Button.SetEnabled(true)
								}
							},
						},

						d.PushButton{
							AssignTo: &stopServerButton,
							Text:     "Stop server",
							Enabled:  false,
							OnClicked: func() {
								if stopServerCallback() {
									startServerButton.Button.SetEnabled(true)
									stopServerButton.Button.SetEnabled(false)
								}
							},
						},
					},
				},

				d.Label{
					Text: "Coils:",
				},

				d.TableView{
					AssignTo:         &coilsView,
					Model:            coilsModel,
					AlternatingRowBG: true,
					ColumnsOrderable: true,
					MaxSize:          d.Size{Height: 125},
					Columns: []d.TableViewColumn{
						{
							Title:      "#",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("%d", value) },
						},
						{
							Title:      "Address",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("0x%X", value) },
						},
						{
							Title:      "Coil",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("%v", value) },
						},
					},
				},

				d.Label{
					Text: "Registers:",
				},

				d.TableView{
					AssignTo:         &registersView,
					Model:            registersModel,
					AlternatingRowBG: true,
					ColumnsOrderable: true,
					MaxSize:          d.Size{Height: 125},
					Columns: []d.TableViewColumn{
						{
							Title:      "#",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("%d", value) },
						},
						{
							Title:      "Address",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("0x%X", value) },
						},
						{
							Title:      "Hex",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("0x%X", value) },
						},
						{
							Title:      "Decimal",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("%d", value) },
						},
					},
				},

				d.Label{
					Text: "Input registers:",
				},

				d.TableView{
					AssignTo:         &inputRegisterView,
					Model:            inputRegistersModel,
					AlternatingRowBG: true,
					ColumnsOrderable: true,
					MaxSize:          d.Size{Height: 125},
					Columns: []d.TableViewColumn{
						{
							Title:      "#",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("%d", value) },
						},
						{
							Title:      "Address",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("0x%X", value) },
						},
						{
							Title:      "Hex",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("0x%X", value) },
						},
						{
							Title:      "Decimal",
							FormatFunc: func(value interface{}) string { return fmt.Sprintf("%d", value) },
						},
					},
				},

				d.TextEdit{
					Name: "log view",
				},
			},
		},
	}
}
