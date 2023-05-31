package main

import (
	"github.com/lxn/walk"
	d "github.com/lxn/walk/declarative"
)

type MainModel interface {
	Connect(transport, address, port string) error
	Reconnect() error
	Disconnect() error

	ReadCoils(addr uint16, cnt int) ([]bool, error)
	ReadDiscreteInputs(addr uint16, cnt int) ([]bool, error)
	ReadHoldingRegisters(addr uint16, cnt int) ([]uint16, error)
	ReadInputRegisters(addr uint16, cnt int) ([]uint16, error)
	WriteSingleCoil(addr uint16, value bool) error
	WriteSingleRegister(addr uint16, value uint16) error
	WriteMultipleRegisters(addr uint16, values []uint16) error
	WriteMultipleCoils(addr uint16, values []bool) error
}

type MainController struct {
	model MainModel

	connParamsSaved bool
	connEstablished bool

	window                       *walk.MainWindow
	connectButton                *walk.PushButton
	reconnectButton              *walk.PushButton
	disconnectButton             *walk.PushButton
	transportEdit                *walk.TextEdit
	addressEdit                  *walk.TextEdit
	portEdit                     *walk.TextEdit
	readCoilsButton              *walk.PushButton
	readDiscreteInputsButton     *walk.PushButton
	readHoldingRegistersButton   *walk.PushButton
	readInputRegistersButton     *walk.PushButton
	writeSingleCoilButton        *walk.PushButton
	writeSingleRegisterButton    *walk.PushButton
	writeMultipleRegistersButton *walk.PushButton
	writeMultipleCoilsButton     *walk.PushButton
	errEdit                      *walk.TextEdit
}

func (c *MainController) Connect() {
	defer c.resetButtons()

	tansport := c.transportEdit.Text()
	address := c.addressEdit.Text()
	port := c.portEdit.Text()
	c.connParamsSaved = true

	if err := c.model.Connect(tansport, address, port); err != nil {
		c.setError(err)
		return
	}

	c.connEstablished = true
	c.clearError()
}

func (c *MainController) Reconnect() {
	defer c.resetButtons()

	if err := c.model.Reconnect(); err != nil {
		c.connEstablished = false
		c.setError(err)
		return
	}

	c.connEstablished = true
	c.clearError()
}

func (c *MainController) Disconnect() {
	defer c.resetButtons()

	if err := c.model.Disconnect(); err != nil {
		c.setError(err)
		return
	}

	c.connEstablished = false
	c.clearError()
}

func (c *MainController) ReadCoils() {
	c.clearError()
	DialogView(
		c.window,
		&DialogModelImpl{c.model},
		DialogTypeReadCoils,
	)()
}

func (c *MainController) ReadDiscreteInputs() {
	c.clearError()
	DialogView(
		c.window,
		&DialogModelImpl{c.model},
		DialogTypeReadDiscreteInputs,
	)()
}

func (c *MainController) ReadHoldingRegisters() {
	c.clearError()
	DialogView(
		c.window,
		&DialogModelImpl{c.model},
		DialogTypeReadHoldingRegisters,
	)()
}

func (c *MainController) ReadInputRegisters() {
	c.clearError()
	DialogView(
		c.window,
		&DialogModelImpl{c.model},
		DialogTypeReadInputRegisters,
	)()
}

func (c *MainController) WriteSingleCoil() {
	c.clearError()
	DialogView(
		c.window,
		&DialogModelImpl{c.model},
		DialogTypeWriteSingleCoil,
	)()
}

func (c *MainController) WriteSingleRegister() {
	c.clearError()
	DialogView(
		c.window,
		&DialogModelImpl{c.model},
		DialogTypeWriteSingleRegister,
	)()
}

func (c *MainController) WriteMultipleRegisters() {
	c.clearError()
	DialogView(
		c.window,
		&DialogModelImpl{c.model},
		DialogTypeWriteMultipleRegisters,
	)()
}

func (c *MainController) WriteMultipleCoils() {
	c.clearError()
	DialogView(
		c.window,
		&DialogModelImpl{c.model},
		DialogTypeWriteMultipleCoils,
	)()
}

func (c *MainController) resetConnectButton() {
	if c.connEstablished {
		c.connectButton.SetEnabled(false)
		return
	}

	c.connectButton.SetEnabled(true)
}

func (c *MainController) resetReconnectButton() {
	if c.connParamsSaved {
		c.reconnectButton.SetEnabled(true)
		return
	}

	c.connectButton.SetEnabled(false)
}

func (c *MainController) resetDisconnectButton() {
	if c.connEstablished {
		c.disconnectButton.SetEnabled(true)
		return
	}

	c.disconnectButton.SetEnabled(false)
}

func (c *MainController) resetFunctionButtons() {
	buttons := []*walk.PushButton{
		c.readCoilsButton,
		c.readDiscreteInputsButton,
		c.readHoldingRegistersButton,
		c.readInputRegistersButton,
		c.writeSingleCoilButton,
		c.writeSingleRegisterButton,
		c.writeMultipleRegistersButton,
		c.writeMultipleCoilsButton,
	}

	for _, b := range buttons {
		if c.connEstablished {
			b.SetEnabled(true)
			continue
		}

		b.SetEnabled(false)
	}
}

func (c *MainController) resetButtons() {
	c.resetConnectButton()
	c.resetReconnectButton()
	c.resetDisconnectButton()
	c.resetFunctionButtons()
}

func (c *MainController) setError(err error) {
	c.errEdit.SetText(err.Error())
}

func (c *MainController) clearError() {
	c.errEdit.SetText("")
}

func MainView(model MainModel) func() {
	controller := &MainController{
		model: model,
	}

	return func() {
		d.MainWindow{
			AssignTo: &controller.window,
			Title:    "Modbus client (master)",
			Size:     d.Size{Width: 320, Height: 300},
			Layout:   d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
			Children: []d.Widget{
				d.GroupBox{
					Title:  "Modbus server address",
					Layout: d.VBox{Margins: d.Margins{Left: 10, Right: 10, Top: 10, Bottom: 10}},
					Children: []d.Widget{

						d.HSplitter{
							Children: []d.Widget{
								d.VSplitter{
									Children: []d.Widget{
										d.Label{Text: "Transport:"},
										d.VSpacer{},
										d.Label{Text: "Address:"},
										d.VSpacer{},
										d.Label{Text: "Port:"},
									},
								},

								d.HSpacer{},

								d.VSplitter{
									Children: []d.Widget{
										d.TextEdit{AssignTo: &controller.transportEdit, Text: "tcp"},
										d.TextEdit{AssignTo: &controller.addressEdit, Text: "localhost"},
										d.TextEdit{AssignTo: &controller.portEdit, Text: "5502"},
									},
								},
							},
						},
					},
				},

				d.HSplitter{
					Children: []d.Widget{
						d.PushButton{
							AssignTo:  &controller.connectButton,
							Text:      "Connect",
							OnClicked: controller.Connect,
						},

						d.PushButton{
							AssignTo:  &controller.reconnectButton,
							Text:      "Reconnect",
							OnClicked: controller.Reconnect,
							Enabled:   false,
						},

						d.PushButton{
							AssignTo:  &controller.disconnectButton,
							Text:      "Disconnect",
							OnClicked: controller.Disconnect,
							Enabled:   false,
						},
					},
				},

				d.VSeparator{},

				d.Composite{
					Layout: d.VBox{Margins: d.Margins{Left: 50, Right: 50}},
					Children: []d.Widget{
						d.PushButton{
							AssignTo:  &controller.readCoilsButton,
							Text:      "0x01 Read coils",
							OnClicked: controller.ReadCoils,
							Enabled:   false,
						},

						d.PushButton{
							AssignTo:  &controller.readDiscreteInputsButton,
							Text:      "0x02 Read discrete inputs",
							OnClicked: controller.ReadDiscreteInputs,
							Enabled:   false,
						},

						d.PushButton{
							AssignTo:  &controller.readHoldingRegistersButton,
							Text:      "0x03 Read holding registers",
							OnClicked: controller.ReadHoldingRegisters,
							Enabled:   false,
						},

						d.PushButton{
							AssignTo:  &controller.readInputRegistersButton,
							Text:      "0x04 Read input registers",
							OnClicked: controller.ReadInputRegisters,
							Enabled:   false,
						},

						d.PushButton{
							AssignTo:  &controller.writeSingleCoilButton,
							Text:      "0x05 Write single coil",
							OnClicked: controller.WriteSingleCoil,
							Enabled:   false,
						},

						d.PushButton{
							AssignTo:  &controller.writeSingleRegisterButton,
							Text:      "0x06 Write single register",
							OnClicked: controller.WriteSingleRegister,
							Enabled:   false,
						},

						d.PushButton{
							AssignTo:  &controller.writeMultipleRegistersButton,
							Text:      "0x10 Write mutliple registers",
							OnClicked: controller.WriteMultipleRegisters,
							Enabled:   false,
						},

						d.PushButton{
							AssignTo:  &controller.writeMultipleCoilsButton,
							Text:      "0x0F Write multiple coils",
							OnClicked: controller.WriteMultipleCoils,
							Enabled:   false,
						},
					},
				},

				d.Label{Text: "Errors:"},
				d.TextEdit{
					AssignTo:  &controller.errEdit,
					MinSize:   d.Size{Height: 50},
					TextColor: walk.RGB(255, 0, 0),
					ReadOnly:  true,
				},
			},
		}.Run()
	}
}

type DialogModelImpl struct {
	MainModel MainModel
}

func (d *DialogModelImpl) ReadCoils(addr uint16, cnt int) ([]bool, error) {
	return d.MainModel.ReadCoils(addr, cnt)
}

func (d *DialogModelImpl) ReadDiscreteInputs(addr uint16, cnt int) ([]bool, error) {
	return d.MainModel.ReadDiscreteInputs(addr, cnt)
}

func (d *DialogModelImpl) ReadHoldingRegisters(addr uint16, cnt int) ([]uint16, error) {
	return d.MainModel.ReadHoldingRegisters(addr, cnt)
}

func (d *DialogModelImpl) ReadInputRegisters(addr uint16, cnt int) ([]uint16, error) {
	return d.MainModel.ReadInputRegisters(addr, cnt)
}

func (d *DialogModelImpl) WriteSingleCoil(addr uint16, value bool) error {
	return d.MainModel.WriteSingleCoil(addr, value)
}

func (d *DialogModelImpl) WriteSingleRegister(addr uint16, value uint16) error {
	return d.MainModel.WriteSingleRegister(addr, value)
}

func (d *DialogModelImpl) WriteMultipleRegisters(addr uint16, values []uint16) error {
	return d.MainModel.WriteMultipleRegisters(addr, values)
}

func (d *DialogModelImpl) WriteMultipleCoils(addr uint16, values []bool) error {
	return d.MainModel.WriteMultipleCoils(addr, values)
}
