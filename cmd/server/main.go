package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/simonvetter/modbus"

	"github.com/aveplen/mirea-modbus/internal/logging"
)

var (
	ErrNoSuchCoil          = errors.New("no such coil")
	ErrNoSuchInputRegister = errors.New("no such input register")
	ErrNoSuchRegister      = errors.New("no such register")
)

// TODO: validate which set of registers should be subscribable and which is not
type SimpleHandler struct {
	coils    map[uint16]bool
	coilSubs []func(update *CoilUpdate)

	inputRegisters    map[uint16]uint16
	inputRegisterSubs []func(update *InputRegisterUpdate)

	registers    map[uint16]uint16
	registerSubs []func(update *InputRegisterUpdate)

	logger *logging.Logger
}

func (h *SimpleHandler) SubscribeToCoils(sub func(update *CoilUpdate)) {
	h.coilSubs = append(h.coilSubs, sub)
}

func (h *SimpleHandler) SubscribeToInputRegisters(sub func(update *InputRegisterUpdate)) {
	h.inputRegisterSubs = append(h.inputRegisterSubs, sub)
}

func (h *SimpleHandler) SubscribeToRegisters(sub func(update *InputRegisterUpdate)) {
	h.registerSubs = append(h.registerSubs, sub)
}

func (h *SimpleHandler) getCoil(addr uint16) (bool, error) {
	value, ok := h.coils[addr]

	if !ok {
		return false, ErrNoSuchCoil
	}

	return value, nil
}

func (h *SimpleHandler) putCoil(addr uint16, value bool) error {
	current, ok := h.coils[addr]

	if !ok {
		return ErrNoSuchCoil
	}

	if current == value {
		return nil
	}

	update := &CoilUpdate{addr: addr, value: value}
	for _, sub := range h.coilSubs {
		sub(update)
	}

	return nil
}

func (h *SimpleHandler) getInputRegister(addr uint16) (uint16, error) {
	value, ok := h.inputRegisters[addr]

	if !ok {
		return 0, ErrNoSuchInputRegister
	}

	return value, nil
}

func (h *SimpleHandler) putInputRegister(addr uint16, value uint16) error {
	current, ok := h.inputRegisters[addr]

	if !ok {
		return ErrNoSuchInputRegister
	}

	if current == value {
		return nil
	}

	update := &InputRegisterUpdate{addr: addr, value: value}
	for _, sub := range h.inputRegisterSubs {
		sub(update)
	}

	return nil
}

func (h *SimpleHandler) getRegister(addr uint16) (uint16, error) {
	value, ok := h.registers[addr]

	if !ok {
		return 0, ErrNoSuchInputRegister
	}

	return value, nil
}

func (h *SimpleHandler) putRegister(addr uint16, value uint16) error {
	current, ok := h.registers[addr]

	if !ok {
		return ErrNoSuchInputRegister
	}

	if current == value {
		return nil
	}

	update := &InputRegisterUpdate{addr: addr, value: value}
	for _, sub := range h.registerSubs {
		sub(update)
	}

	return nil
}

func NewSimpleHandler(
	logger *logging.Logger,
	coils map[uint16]bool,
	inputRegisters map[uint16]uint16,
	registers map[uint16]uint16,
) *SimpleHandler {

	handler := &SimpleHandler{
		logger:         logger,
		coils:          coils,
		inputRegisters: inputRegisters,
		registers:      registers,
	}

	return handler
}

// HandleCoils handles the read coils (0x01), write single coil (0x05)
// and write multiple coils (0x0F)
// - res:	coil values (only for reads)
// - err:	either nil if no error occurred, a modbus error
func (h *SimpleHandler) HandleCoils(req *modbus.CoilsRequest) ([]bool, error) {
	if req.UnitId != 1 {
		h.logger.Errorf("HandleCoils accessed with wrong UnitId: %d", req.UnitId)
		return []bool{false}, modbus.ErrIllegalFunction
	}

	if req.IsWrite && req.Quantity == 1 {
		h.logger.Debugf(
			"Function 0x05 (write single coil) accessed with Addr: %d, Arguments: %v",
			req.Addr, req.Args,
		)

		if err := h.putCoil(req.Addr, req.Args[0]); err != nil {
			h.logger.Error(err)
			return []bool{false}, modbus.ErrIllegalDataAddress
		}

		return []bool{true}, nil
	}

	if req.IsWrite {
		h.logger.Debugf(
			"Function 0x0F (write multiple coils) accessed with Addr: %d, Quantity: %d, Arguments: %v",
			req.Addr, req.Quantity, req.Args,
		)

		for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
			argNo := addr - req.Addr
			if err := h.putCoil(req.Addr, req.Args[argNo]); err != nil {
				h.logger.Error(err)
				return []bool{false}, modbus.ErrIllegalDataAddress
			}
		}

		return []bool{true}, nil
	}

	h.logger.Debugf(
		"Funtion 0x01 (read coils) accessed with Addr: %d, Quantity: %d",
		req.Addr, req.Quantity,
	)

	res := make([]bool, 0, req.Quantity)
	for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
		coil, err := h.getCoil(addr)

		if err != nil {
			h.logger.Error(err)
			return []bool{false}, modbus.ErrIllegalDataAddress
		}

		res = append(res, coil)
	}

	return res, nil
}

// HandleDiscreteInputs handles the read discrete inputs (0x02)
// - res: discrete input values
// - err:	either nil if no error occurred, a modbus error
func (h *SimpleHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) ([]bool, error) {
	if req.UnitId != 1 {
		h.logger.Errorf("HandleDiscreteInputs accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}

	if req.Quantity != 1 {
		h.logger.Debugf("HandleDiscreteInputs accessed with unsupported Quantity: %d", req.Quantity)
		return nil, modbus.ErrIllegalDataAddress
	}

	h.logger.Infof("Function 0x02 (read discrete inputs) at Addr: %d with Quantity: %d", req.Addr, req.Quantity)

	if req.Addr == 10071 {
		return []bool{false}, nil
	}

	if req.Addr == 14012 {
		d10 := rand.Int() % 10

		if d10 == 0 {
			return nil, modbus.ErrServerDeviceFailure
		}

		if d10 <= 5 {
			return []bool{true}, nil
		}

		return []bool{false}, nil
	}

	return nil, modbus.ErrIllegalDataAddress
}

// HandleHoldingRegisters handles the read holding registers (0x03),
// write single register (0x06) and write multiple registers (0x10).
// A HoldingRegistersRequest object is passed to the handler (see above).
// - res:	register values
// - err:	either nil if no error occurred, a modbus error
func (h *SimpleHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) ([]uint16, error) {
	if req.UnitId != 1 {
		h.logger.Errorf("HandleHoldingRegisters accessed with wrong UnitId: %d", req.UnitId)
		return []uint16{0}, modbus.ErrIllegalFunction
	}

	if req.IsWrite && req.Quantity == 1 {
		h.logger.Debugf(
			"Function 0x06 (write single register) accessed with Addr: %d, Quantity: %d, Arguments: %v",
			req.Addr, req.Quantity, req.Args,
		)

		if err := h.putInputRegister(req.Addr, req.Args[0]); err != nil {
			h.logger.Errorf("Failed attempt to write single register at addr %d", req.Addr)
			return []uint16{0}, modbus.ErrIllegalDataAddress
		}

		return []uint16{0}, nil
	}

	if req.IsWrite {
		h.logger.Debugf(
			"Function 0x10 (write mutliple registers) accessed with Addr: %d, Quantity: %d, Arguments: %v",
			req.Addr, req.Quantity, req.Args,
		)

		for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
			argNo := addr - req.Addr
			if err := h.putInputRegister(req.Addr, req.Args[argNo]); err != nil {
				h.logger.Errorf("Failed attempt to write mutliple registers at addr %d", addr)
				return []uint16{0}, modbus.ErrIllegalDataAddress
			}
		}

		return []uint16{0}, nil
	}

	log.Printf("Function 0x03 (read holding registers) accessed with Addr: %d, Quantity: %d", req.Addr, req.Quantity)
	res := make([]uint16, 0, req.Quantity)
	for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
		value, err := h.getInputRegister(addr)
		if err != nil {
			h.logger.Errorf("Failed attempt to read register at addr %d", addr)
			return []uint16{0}, modbus.ErrIllegalDataAddress
		}

		res = append(res, value)
	}

	return res, nil
}

// HandleInputRegisters handles the read input registers (0x04)
// Note that input registers are always read-only as per the modbus spec.
// - res:	register values
// - err:	either nil if no error occurred, a modbus error
func (h *SimpleHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) ([]uint16, error) {
	if req.UnitId != 1 {
		h.logger.Errorf("HandleInputRegisters accessed with wrong UnitId: %d", req.UnitId)
		return []uint16{0}, modbus.ErrIllegalFunction
	}

	h.logger.Debugf("Function 0x04 (read input registers) at Addr: %d with Quantity: %d", req.Addr, req.Quantity)

	if req.Quantity != 2 {
		h.logger.Debugf("Unsupported operation: reading %d input registers", req.Quantity)
		return []uint16{0}, modbus.ErrIllegalFunction
	}

	for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
		if _, err := h.getRegister(addr); err != nil {
			h.logger.Errorf("Failed attempt to read input register at addr %d", addr)
			return []uint16{0}, modbus.ErrIllegalDataAddress
		}
	}

	bytes := []byte{1, 0, 3, 2}
	pairs := make([][]byte, 0, len(bytes)/2)
	for i, b := range bytes {
		pairIndex := i / 2

		if i%2 == 0 {
			pairs = append(pairs, make([]byte, 0, 2))
		}

		pairs[pairIndex] = append(pairs[pairIndex], b)
	}

	res := make([]uint16, 0, len(pairs))
	for _, pair := range pairs {
		res = append(res, binary.BigEndian.Uint16(pair))
	}

	return res, nil
}

type ChanSubscriber struct {
	messages chan string
}

func NewChanSubscriber(messages chan string) *ChanSubscriber {
	listSubscriber := &ChanSubscriber{
		messages: messages,
	}
	return listSubscriber
}

func (l *ChanSubscriber) Consume(message string) {
	l.messages <- message
}

type CoilUpdate struct {
	addr  uint16
	value bool
}

type CoilsView struct {
	coils map[uint16]bool
	table *widgets.Table
}

func NewCoilsView(seed map[uint16]bool) *CoilsView {
	table := widgets.NewTable()
	table.Block.Title = "Coils"
	table.Rows = coilsToTable(seed)
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 0, 80, 5)

	view := &CoilsView{
		coils: seed,
		table: table,
	}

	return view
}

func coilsToTable(coils map[uint16]bool) [][]string {
	order := make([]uint16, 0, len(coils))
	for addr := range coils {
		order = append(order, addr)
	}

	sort.Slice(order, func(i, j int) bool {
		return order[i] < order[j]
	})

	header := make([]string, 0, len(coils))
	values := make([]string, 0, len(coils))
	for _, addr := range order {
		header = append(header, fmt.Sprint(addr))
		values = append(values, fmt.Sprint(coils[uint16(addr)]))
	}

	return [][]string{header, values}
}

func (c *CoilsView) HandleUpdate(update *CoilUpdate) {
	c.coils[update.addr] = update.value
	c.table.Rows = coilsToTable(c.coils)
}

func (c *CoilsView) GetComponent() *widgets.Table {
	return c.table
}

type InputRegisterUpdate struct {
	addr  uint16
	value uint16
}

type InputRegistersView struct {
	registers map[uint16]uint16
	table     *widgets.Table
}

func NewInputRegistersView(seed map[uint16]uint16) *InputRegistersView {
	table := widgets.NewTable()
	table.Block.Title = "Input registers"
	table.Rows = inputRegistersToTable(seed)
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 6, 80, 11)

	view := &InputRegistersView{
		table:     table,
		registers: seed,
	}

	return view
}

func inputRegistersToTable(inputRegisters map[uint16]uint16) [][]string {
	order := make([]uint16, 0, len(inputRegisters))
	for addr := range inputRegisters {
		order = append(order, addr)
	}

	sort.Slice(order, func(i, j int) bool {
		return order[i] < order[j]
	})

	header := make([]string, 0, len(inputRegisters))
	values := make([]string, 0, len(inputRegisters))
	for _, addr := range order {
		header = append(header, fmt.Sprint(addr))
		values = append(values, fmt.Sprint(inputRegisters[uint16(addr)]))
	}

	return [][]string{header, values}
}

func (r *InputRegistersView) HandleUpdate(update *InputRegisterUpdate) {
	r.registers[update.addr] = update.value
	r.table.Rows = inputRegistersToTable(r.registers)
}

func (r *InputRegistersView) GetComponent() *widgets.Table {
	return r.table
}

type LogView struct {
	messages []string
	list     *widgets.List
}

func NewLogView() *LogView {
	list := widgets.NewList()
	list.Title = "Log"
	list.Rows = []string{}
	list.WrapText = false
	list.SetRect(0, 12, 80, 30)

	view := &LogView{
		messages: make([]string, 0, 25),
		list:     list,
	}

	return view
}

func (l *LogView) HandleUpdate(message string) {
	if len(l.messages) == 0 {
		l.messages = append(l.messages, message)
		l.list.Rows = logsToRows(l.messages)
		return
	}

	l.messages = append(l.messages[1:], message)
	l.list.Rows = logsToRows(l.messages)
}

func logsToRows(logs []string) []string {
	return logs
}

func (l *LogView) GetComponent() *widgets.List {
	return l.list
}

func LogList() *widgets.List {
	l := widgets.NewList()
	l.Title = "Log"
	l.Rows = make([]string, 0, 10)
	l.WrapText = false
	l.SetRect(0, 12, 80, 30)
	return l
}

var (
	coils = map[uint16]bool{
		21560: false,
		21561: false,
		21565: false,
		21566: false,
	}

	// TODO: check which set of registers should be modifiable and which is not

	// TODO validate this
	registers = map[uint16]uint16{
		30022: 0,
		30023: 0,
		30024: 0,
		30025: 0,
		30026: 0,
		30027: 0,
		30028: 0,
	}

	// TODO validate this
	inputRegisters = map[uint16]uint16{
		44883: 0,
		44884: 0,
		44885: 0,
		44886: 0,
		44887: 0,
		44888: 0,
		44889: 0,
		44890: 0,
	}
)

func main() {
	// logger
	logUpdates := make(chan string, 100)
	chanSubscriber := NewChanSubscriber(logUpdates)
	loggerFactory := logging.NewLoggerFactory(
		logging.NewLoggerConfig(
			logging.WithTemplate("%v %v"),
			logging.WithSubscribers(chanSubscriber),
		),
	)

	// handler
	handler := NewSimpleHandler(
		loggerFactory.GetLogger("SimpleHandler"),
		coils,
		inputRegisters,
		registers,
	)

	coilUpdates := make(chan *CoilUpdate)
	handler.SubscribeToCoils(func(update *CoilUpdate) {
		coilUpdates <- update
	})

	inputRegisterUpdates := make(chan *InputRegisterUpdate)
	handler.SubscribeToInputRegisters(func(update *InputRegisterUpdate) {
		inputRegisterUpdates <- update
	})

	// server
	server, err := modbus.NewServer(
		&modbus.ServerConfiguration{
			URL:        "tcp://localhost:5502",
			Timeout:    30 * time.Second,
			MaxClients: 5,
		},
		handler,
	)

	if err != nil {
		panic(fmt.Errorf("failed to create server: %w", err))
	}

	go func() {
		if err := server.Start(); err != nil {
			panic(fmt.Errorf("failed to start server: %w", err))
		}
		log.Println("server started")
	}()

	// ui
	coilsView := NewCoilsView(coils)
	inputRegistersView := NewInputRegistersView(inputRegisters)
	logView := NewLogView()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	uiEvents := ui.PollEvents()

	coilsViewCompoent := coilsView.GetComponent()
	inputRegistersViewCompoennt := inputRegistersView.GetComponent()
	logViewCompoenent := logView.GetComponent()

	ui.Render(
		coilsViewCompoent,
		inputRegistersViewCompoennt,
		logViewCompoenent,
	)

	ticker := time.NewTicker(250 * time.Millisecond)
	mainLogger := loggerFactory.GetLogger("Main Logger")

	for {
		select {
		case t := <-ticker.C:
			{
				mainLogger.Info(t.String())
				ui.Render(logViewCompoenent)
			}

		case coilsUpdate := <-coilUpdates:
			{
				coilsView.HandleUpdate(coilsUpdate)
				ui.Render(coilsViewCompoent)
			}

		case inputRegisterUpdate := <-inputRegisterUpdates:
			{
				inputRegistersView.HandleUpdate(inputRegisterUpdate)
				ui.Render(inputRegistersViewCompoennt)
			}

		case logUpdate := <-logUpdates:
			{
				logView.HandleUpdate(logUpdate)
				ui.Render(logViewCompoenent)
			}

		case event := <-uiEvents:
			{
				if event.ID == "q" || event.ID == "<C-c>" {
					func() { // for defer cancel()
						ctx := context.Background()
						timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
						defer cancel()

						stopped := make(chan error)
						go func() {
							stopped <- server.Stop()
						}()

						select {
						case <-timeout.Done():
							{
								panic(errors.New("could not stop server in time"))
							}

						case err := <-stopped:
							{
								if err != nil {
									panic(fmt.Errorf("server stopped with error: %w", err))
								}

								log.Println("server stopped")
								os.Exit(0)
							}
						}
					}()
				}
			}
		}
	}
}
