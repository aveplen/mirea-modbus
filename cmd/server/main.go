package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aveplen/mirea-modbus/internal/util"
	"github.com/simonvetter/modbus"
)

type SimpleHandler struct {
	ctx  context.Context
	lock *util.TriableMutex

	startedAt         int64
	coils             map[uint16]bool
	readonlyRegisters map[uint16]uint16
	regs              map[uint16]uint16

	uptime       int64
	tickerCancel context.CancelFunc
}

func NewSimpleHandler(ctx context.Context) *SimpleHandler {
	coils := make(map[uint16]bool)
	coils[21560] = false
	coils[21561] = false
	coils[21565] = false
	coils[21566] = false

	readonlyRegisters := make(map[uint16]uint16)
	readonlyRegisters[30022] = 0
	readonlyRegisters[30023] = 0
	readonlyRegisters[30023] = 0
	readonlyRegisters[30024] = 0
	readonlyRegisters[30025] = 0
	readonlyRegisters[30026] = 0
	readonlyRegisters[30027] = 0
	readonlyRegisters[30028] = 0

	regs := make(map[uint16]uint16)
	regs[44883] = 0
	regs[44884] = 0
	regs[44885] = 0
	regs[44886] = 0
	regs[44887] = 0
	regs[44888] = 0
	regs[44889] = 0
	regs[44890] = 0

	handler := &SimpleHandler{
		ctx:               ctx,
		lock:              util.NewTryableMutex(),
		coils:             coils,
		readonlyRegisters: readonlyRegisters,
		regs:              regs,
	}
	return handler
}

func (h *SimpleHandler) Start() {
	h.startedAt = time.Now().Unix()
	h.startTicker()
}

func (h *SimpleHandler) Stop() {
	h.stopTicker()
}

func (h *SimpleHandler) startTicker() {
	ctxCancel, cancelCancel := context.WithCancel(h.ctx)
	h.tickerCancel = cancelCancel

	log.Println("starting ticker ...")

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ctxCancel.Done():
				return

			case tick := <-ticker.C:
				go func() {
					acquired := <-h.lock.LockTimeout(time.Second)
					if !acquired {
						return
					}

					defer h.lock.Unlock()

					sinceStarted := tick.Unix() - h.startedAt
					if sinceStarted > h.uptime {
						h.uptime = sinceStarted
						// log.Printf("current handler uptime: %d\n", h.uptime)
						return
					}
				}()
			}
		}
	}()
}

func (h *SimpleHandler) stopTicker() {
	h.tickerCancel()
}

func (h *SimpleHandler) currentUptime() int64 {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.uptime
}

// HandleCoils handles the read coils (0x01), write single coil (0x05)
// and write multiple coils (0x0F)
// - res:	coil values (only for reads)
// - err:	either nil if no error occurred, a modbus error
func (h *SimpleHandler) HandleCoils(req *modbus.CoilsRequest) ([]bool, error) {
	if req.UnitId != 1 {
		log.Printf("HandleCoils accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}

	if req.IsWrite && req.Quantity == 1 {
		log.Printf(
			"Function 0x05 (write single coil) accessed with Addr: %d, Arguments: %v",
			req.Addr, req.Args,
		)

		if _, ok := h.coils[req.Addr]; !ok {
			return nil, modbus.ErrIllegalDataAddress
		}

		h.coils[req.Addr] = req.Args[0]
		return nil, nil
	}

	if req.IsWrite {
		log.Printf(
			"Function 0x0F (write multiple coils) accessed with Addr: %d, Quantity: %d, Arguments: %v",
			req.Addr, req.Quantity, req.Args,
		)

		for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
			if _, ok := h.coils[addr]; !ok {
				return nil, modbus.ErrIllegalDataAddress
			}
			argNo := addr - req.Addr
			h.coils[req.Addr] = req.Args[argNo]
		}
	}

	log.Printf(
		"Funtion 0x01 (read coils) accessed with Addr: %d, Quantity: %d",
		req.Addr, req.Quantity,
	)

	res := make([]bool, 0, req.Quantity)
	for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
		coil, ok := h.coils[addr]
		if !ok {
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
		log.Printf("HandleDiscreteInputs accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}

	if req.Quantity != 1 {
		log.Printf("HandleDiscreteInputs accessed with unsupported Quantity: %d", req.Quantity)
		return nil, modbus.ErrIllegalDataAddress
	}

	log.Printf("Function 0x02 (read discrete inputs) at Addr: %d with Quantity: %d", req.Addr, req.Quantity)

	if req.Addr == 10071 {
		log.Printf("Retutning current server uptime mod 2")
		uptime := h.currentUptime()
		uptimeMod2 := uptime%2 == 1
		log.Printf("Current uptime: %d, mod 2: %v", uptime, uptimeMod2)
		return []bool{uptimeMod2}, nil
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
		log.Printf("HandleHoldingRegisters accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}

	if req.IsWrite && req.Quantity == 1 {
		log.Printf(
			"Function 0x06 (write single register) accessed with Addr: %d, Quantity: %d, Arguments: %v",
			req.Addr, req.Quantity, req.Args,
		)

		if _, ok := h.regs[req.Addr]; !ok {
			log.Printf("Failed attempt to write single register at addr %d", req.Addr)
			return nil, modbus.ErrIllegalDataAddress
		}

		h.regs[req.Addr] = req.Args[0]
		return nil, nil
	}

	if req.IsWrite {
		log.Printf(
			"Function 0x10 (write mutliple registers) accessed with Addr: %d, Quantity: %d, Arguments: %v",
			req.Addr, req.Quantity, req.Args,
		)

		for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
			if _, ok := h.regs[addr]; !ok {
				log.Printf("Failed attempt to write mutliple registers at addr %d", addr)
				return nil, modbus.ErrIllegalDataAddress
			}

			argNo := addr - req.Addr
			h.regs[addr] = req.Args[argNo]
		}

		return nil, nil
	}

	log.Printf("Function 0x03 (read holding registers) accessed with Addr: %d, Quantity: %d", req.Addr, req.Quantity)
	res := make([]uint16, 0, req.Quantity)
	for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
		value, ok := h.regs[addr]

		if !ok {
			log.Printf("Failed attempt to read register at addr %d", addr)
			return nil, modbus.ErrIllegalDataAddress
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
		log.Printf("HandleInputRegisters accessed with wrong UnitId: %d", req.UnitId)
		return nil, modbus.ErrIllegalFunction
	}

	log.Printf("Function 0x04 (read input registers) at Addr: %d with Quantity: %d", req.Addr, req.Quantity)

	if req.Quantity != 2 {
		log.Printf("Unsupported operation: reading %d input registers", req.Quantity)
		return nil, modbus.ErrIllegalFunction
	}

	for addr := req.Addr; addr < req.Addr+req.Quantity; addr++ {
		if _, ok := h.readonlyRegisters[addr]; !ok {
			return nil, modbus.ErrIllegalDataAddress
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

func main() {
	ctx := context.Background()

	handler := NewSimpleHandler(ctx)
	handler.Start()

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

	quit := make(chan os.Signal, 5)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stopped := make(chan error)
	go func() {
		stopped <- server.Stop()
	}()

	select {
	case <-timeout.Done():
		panic(errors.New("could not stop server in time"))

	case err := <-stopped:
		if err != nil {
			panic(fmt.Errorf("server stopped with error: %w", err))
		}
		log.Println("server stopped")
		os.Exit(0)
	}
}
