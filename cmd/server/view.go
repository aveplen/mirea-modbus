package main

import (
	"fmt"
	"sort"

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

	sort.Slice(m.items, func(i, j int) bool {
		return m.items[i].Address < m.items[j].Address
	})

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

	sort.Slice(m.items, func(i, j int) bool {
		return m.items[i].Address < m.items[j].Address
	})

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

	StartSimulation()
	StopSimulation()
}

type ViewController struct {
	MainWindow *d.MainWindow
	model      MainModel

	discreteInputsModel   *CoilsModel
	coilsModel            *CoilsModel
	inputRegistersModel   *RegistersModel
	holdingRegistersModel *RegistersModel

	discreteInputsView    *walk.TableView
	coilsView             *walk.TableView
	inputRegisterView     *walk.TableView
	holdingRegistersView  *walk.TableView
	startServerButton     *walk.PushButton
	stopServerButton      *walk.PushButton
	startSimulationButton *walk.PushButton
	stopSimulationButton  *walk.PushButton
	clearLogButton        *walk.PushButton

	AppendLog func(value string)
}

func (v *ViewController) UpdateDiscreteInputs(change CoilChange) {
	for _, item := range v.discreteInputsModel.items {
		if item.Address == change.addr {
			item.Value = change.to
			break
		}
	}

	v.discreteInputsModel.PublishRowsReset()
	v.discreteInputsView.Invalidate()
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

func (v *ViewController) StartSimulation() {
	v.model.StartSimulation()
	v.startSimulationButton.Button.SetEnabled(false)
	v.stopSimulationButton.Button.SetEnabled(true)
}

func (v *ViewController) StopSimulation() {
	v.model.StopSimulation()
	v.startSimulationButton.Button.SetEnabled(true)
	v.stopSimulationButton.Button.SetEnabled(false)
}

func NewView(seed Dump, model MainModel) *ViewController {
	view := &ViewController{
		model: model,

		discreteInputsModel:   NewCoilsModel(seed.DiscreteInputs),
		coilsModel:            NewCoilsModel(seed.Coils),
		inputRegistersModel:   NewRegistersModel(seed.InputRegisters),
		holdingRegistersModel: NewRegistersModel(seed.HoldingRegisters),
	}

	lv := NewLogView()
	view.AppendLog = lv.Append

	view.MainWindow = &d.MainWindow{
		Title:  "Modbus server (slave)",
		Size:   d.Size{Width: 1200, Height: 1000},
		Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
		Children: []d.Widget{
			d.Composite{
				Layout: d.HBox{},
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

					d.PushButton{
						AssignTo:  &view.startSimulationButton,
						Text:      "Start simulation",
						OnClicked: view.StartSimulation,
					},

					d.PushButton{
						AssignTo:  &view.stopSimulationButton,
						Text:      "Stop simulation",
						OnClicked: view.StopSimulation,
						Enabled:   false,
					},
				},
			},

			d.Composite{
				Layout: d.HBox{},
				Children: []d.Widget{
					d.Composite{
						Layout:  d.VBox{},
						MinSize: d.Size{Width: 40},
						Children: []d.Widget{
							d.Label{Text: "Discrete inputs:"},
							d.TableView{
								AssignTo:         &view.discreteInputsView,
								Model:            view.discreteInputsModel,
								AlternatingRowBG: true,
								ColumnsOrderable: true,
								Columns: []d.TableViewColumn{
									{Title: "#", FormatFunc: numberDecimalFormat},
									{Title: "Address (dec)", FormatFunc: numberDecimalFormat},
									{Title: "Address (hex)", FormatFunc: numberHexFormat},
									{Title: "Value", FormatFunc: boolFormat},
								},
							},

							d.Label{Text: "Coils:"},
							d.TableView{
								AssignTo:         &view.coilsView,
								Model:            view.coilsModel,
								AlternatingRowBG: true,
								ColumnsOrderable: true,
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
								Columns: []d.TableViewColumn{
									{Title: "#", FormatFunc: numberDecimalFormat},
									{Title: "Address (dec)", FormatFunc: numberDecimalFormat},
									{Title: "Address (hex)", FormatFunc: numberHexFormat},
									{Title: "Value (dec)", FormatFunc: numberDecimalFormat},
									{Title: "Value (hex)", FormatFunc: numberHexFormat},
								},
							},
						},
					},

					d.Composite{
						Layout: d.VBox{},
						Children: []d.Widget{
							d.Composite{
								Layout: d.HBox{},
								Children: []d.Widget{
									d.Label{Text: "Log view:"},
									d.PushButton{
										AssignTo:  &view.clearLogButton,
										Text:      "Clear",
										OnClicked: lv.ClearLog,
									},
								},
							},

							lv.TextEdit,
						},
					},
				},
			},
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
