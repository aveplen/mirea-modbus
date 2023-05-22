package main

import (
	"github.com/lxn/walk"
	d "github.com/lxn/walk/declarative"
)

type ViewCallbacks struct {
	connect    func(transport, address, port string) error
	reconnect  func() error
	disconnect func() error

	readCoils              func(addr uint16, cnt int) ([]bool, error)
	readDiscreteInputs     func(addr uint16, cnt int) ([]bool, error)
	readHoldingRegisters   func(addr uint16, cnt int) ([]uint16, error)
	readInputRegisters     func(addr uint16, cnt int) ([]uint16, error)
	writeSingleCoil        func(addr uint16, value bool) error
	writeSingleRegister    func(addr uint16, value uint16) error
	writeMultipleRegisters func(addr uint16, values []uint16) error
	writeMultipleCoils     func(addr uint16, values []bool) error
}

type View struct {
	callbacks ViewCallbacks

	connParamsSaved bool
	connEstablished bool

	window *walk.MainWindow

	connectButton    *walk.PushButton
	reconnectButton  *walk.PushButton
	disconnectButton *walk.PushButton

	transportEdit *walk.TextEdit
	addressEdit   *walk.TextEdit
	portEdit      *walk.TextEdit

	readCoilsButton              *walk.PushButton
	readDiscreteInputsButton     *walk.PushButton
	readHoldingRegistersButton   *walk.PushButton
	readInputRegistersButton     *walk.PushButton
	writeSingleCoilButton        *walk.PushButton
	writeSingleRegisterButton    *walk.PushButton
	writeMultipleRegistersButton *walk.PushButton
	writeMultipleCoilsButton     *walk.PushButton

	errEdit *walk.TextEdit
}

func NewView(callbacks ViewCallbacks) *View {
	return &View{
		callbacks: callbacks,
	}
}

func (v *View) resetConnectButton() {
	if v.connEstablished {
		v.connectButton.SetEnabled(false)
		return
	}
	v.connectButton.SetEnabled(true)
}

func (v *View) resetReconnectButton() {
	if v.connParamsSaved {
		v.reconnectButton.SetEnabled(true)
		return
	}
	v.connectButton.SetEnabled(false)
}

func (v *View) resetDisconnectButton() {
	if v.connEstablished {
		v.disconnectButton.SetEnabled(true)
		return
	}
	v.disconnectButton.SetEnabled(false)
}

func (v *View) resetFunctionButtons() {
	buttons := []*walk.PushButton{
		v.readCoilsButton,
		v.readDiscreteInputsButton,
		v.readHoldingRegistersButton,
		v.readInputRegistersButton,
		v.writeSingleCoilButton,
		v.writeSingleRegisterButton,
		v.writeMultipleRegistersButton,
		v.writeMultipleCoilsButton,
	}

	for _, b := range buttons {
		if v.connEstablished {
			b.SetEnabled(true)
			continue
		}

		b.SetEnabled(false)
	}
}

func (v *View) resetButtons() {
	v.resetConnectButton()
	v.resetReconnectButton()
	v.resetDisconnectButton()

	v.resetFunctionButtons()
}

func (v *View) connect() {
	defer v.resetButtons()

	tansport := v.transportEdit.Text()
	address := v.addressEdit.Text()
	port := v.portEdit.Text()
	v.connParamsSaved = true

	if err := v.callbacks.connect(tansport, address, port); err != nil {
		v.errEdit.SetText(err.Error())
		return
	}

	v.connEstablished = true
	v.errEdit.SetText("")
}

func (v *View) reconnect() {
	defer v.resetButtons()

	if err := v.callbacks.reconnect(); err != nil {
		v.errEdit.SetText(err.Error())
		v.connEstablished = false
		return
	}

	v.connEstablished = true
	v.errEdit.SetText("")
}

func (v *View) disconnect() {
	defer v.resetButtons()

	if err := v.callbacks.disconnect(); err != nil {
		v.errEdit.SetText(err.Error())
		return
	}

	v.connEstablished = false
	v.errEdit.SetText("")
}

func (v *View) readCoils() {
	dialog := NewDialog(
		DialogCallbacks{readCoils: v.callbacks.readCoils},
		DialogTypeReadCoils,
	)
	dialog.Render().Run(v.window)
}

func (v *View) readDiscreteInputs() {
	dialog := NewDialog(
		DialogCallbacks{readDiscreteInputs: v.callbacks.readDiscreteInputs},
		DialogTypeReadDiscreteInputs,
	)
	dialog.Render().Run(v.window)
}

func (v *View) readHoldingRegisters() {
	dialog := NewDialog(
		DialogCallbacks{readHoldingRegisters: v.callbacks.readHoldingRegisters},
		DialogTypeReadHoldingRegisters,
	)
	dialog.Render().Run(v.window)
}

func (v *View) readInputRegisters() {
	dialog := NewDialog(
		DialogCallbacks{readInputRegisters: v.callbacks.readInputRegisters},
		DialogTypeReadInputRegisters,
	)
	dialog.Render().Run(v.window)
}

func (v *View) writeSingleCoil() {
	dialog := NewDialog(
		DialogCallbacks{writeSingleCoil: v.callbacks.writeSingleCoil},
		DialogTypeWriteSingleCoil,
	)
	dialog.Render().Run(v.window)
}

func (v *View) writeSingleRegister() {
	dialog := NewDialog(
		DialogCallbacks{writeSingleRegister: v.callbacks.writeSingleRegister},
		DialogTypeWriteSingleRegister,
	)
	dialog.Render().Run(v.window)
}

func (v *View) writeMultipleRegisters() {
	dialog := NewDialog(
		DialogCallbacks{writeMultipleRegisters: v.callbacks.writeMultipleRegisters},
		DialogTypeWriteMultipleRegisters,
	)
	dialog.Render().Run(v.window)
}

func (v *View) writeMultipleCoils() {
	dialog := NewDialog(
		DialogCallbacks{writeMultipleCoils: v.callbacks.writeMultipleCoils},
		DialogTypeWriteMultipleCoils,
	)
	dialog.Render().Run(v.window)
}

func (v *View) Render() *d.MainWindow {
	return &d.MainWindow{
		AssignTo: &v.window,
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
									d.Label{
										Text: "Transport:",
									},

									d.VSpacer{},

									d.Label{
										Text: "Address:",
									},

									d.VSpacer{},

									d.Label{
										Text: "Port:",
									},
								},
							},

							d.HSpacer{},

							d.VSplitter{
								Children: []d.Widget{
									d.TextEdit{
										AssignTo: &v.transportEdit,
										Text:     "tcp",
									},

									d.TextEdit{
										AssignTo: &v.addressEdit,
										Text:     "localhost",
									},

									d.TextEdit{
										AssignTo: &v.portEdit,
										Text:     "5502",
									},
								},
							},
						},
					},
				},
			},

			d.HSplitter{
				Children: []d.Widget{
					d.PushButton{
						AssignTo:  &v.connectButton,
						Text:      "Connect",
						OnClicked: v.connect,
					},

					d.PushButton{
						AssignTo:  &v.reconnectButton,
						Text:      "Reconnect",
						OnClicked: v.reconnect,
					},

					d.PushButton{
						AssignTo:  &v.disconnectButton,
						Text:      "Disconnect",
						OnClicked: v.disconnect,
					},
				},
			},

			d.VSeparator{},

			d.Composite{
				Layout: d.VBox{Margins: d.Margins{Left: 50, Right: 50}},
				Children: []d.Widget{
					d.PushButton{
						AssignTo:  &v.readCoilsButton,
						Text:      "0x01 Read coils",
						Enabled:   false,
						OnClicked: v.readCoils,
					},

					d.PushButton{
						AssignTo:  &v.readDiscreteInputsButton,
						Text:      "0x02 Read discrete inputs",
						Enabled:   false,
						OnClicked: v.readDiscreteInputs,
					},

					d.PushButton{
						AssignTo:  &v.readHoldingRegistersButton,
						Text:      "0x03 Read holding registers",
						Enabled:   false,
						OnClicked: v.readHoldingRegisters,
					},

					d.PushButton{
						AssignTo:  &v.readInputRegistersButton,
						Text:      "0x04 Read input registers",
						Enabled:   false,
						OnClicked: v.readInputRegisters,
					},

					d.PushButton{
						AssignTo:  &v.writeSingleCoilButton,
						Text:      "0x05 Write single coil",
						Enabled:   false,
						OnClicked: v.writeSingleCoil,
					},

					d.PushButton{
						AssignTo:  &v.writeSingleRegisterButton,
						Text:      "0x06 Write single register",
						Enabled:   false,
						OnClicked: v.writeSingleRegister,
					},

					d.PushButton{
						AssignTo:  &v.writeMultipleRegistersButton,
						Text:      "0x10 Write mutliple registers",
						Enabled:   false,
						OnClicked: v.writeMultipleRegisters,
					},

					d.PushButton{
						AssignTo:  &v.writeMultipleCoilsButton,
						Text:      "0x0F Write multiple coils",
						Enabled:   false,
						OnClicked: v.writeMultipleCoils,
					},
				},
			},

			d.Label{
				Text: "Errors:",
			},

			d.TextEdit{
				MinSize:   d.Size{Width: 1, Height: 50},
				ReadOnly:  true,
				TextColor: walk.RGB(255, 0, 0),
				AssignTo:  &v.errEdit,
			},
		},
	}
}
