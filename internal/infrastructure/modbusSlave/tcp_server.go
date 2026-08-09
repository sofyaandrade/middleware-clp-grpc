package modbusSlave

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/jobs"
	"net"
	"sync"
	"time"
)

const (
	modbusFunctionReadCoils            byte = byte(modbusOperationCoilStatus)
	modbusFunctionReadDiscreteInputs   byte = byte(modbusOperationInputStatus)
	modbusFunctionReadHoldingRegisters byte = byte(modbusOperationHoldingRegister)
	modbusFunctionReadInputRegisters   byte = byte(modbusOperationInputRegister)
	modbusExceptionIllegalFunction     byte = 1
	modbusExceptionIllegalDataAddress  byte = 2
	modbusExceptionIllegalDataValue    byte = 3
	modbusExceptionServerFailure       byte = 4
)

type slaveTCPServer struct {
	address string
	slaveID byte
	timeout time.Duration
	clpID   uint
	tags    []*models.Tag

	mu          sync.Mutex
	listener    net.Listener
	connections map[net.Conn]struct{}
	errCh       chan error
	memory      *slaveMemory
}

type slaveMemory struct {
	mu               sync.RWMutex
	coils            map[uint16]bool
	discreteInputs   map[uint16]bool
	holdingRegisters map[uint16]uint16
	inputRegisters   map[uint16]uint16
}

func newSlaveTCPServer(address string, slaveID byte, clpID uint, tags []*models.Tag, timeout time.Duration) *slaveTCPServer {
	return &slaveTCPServer{
		address:     address,
		slaveID:     slaveID,
		timeout:     timeout,
		clpID:       clpID,
		tags:        tags,
		connections: make(map[net.Conn]struct{}),
		errCh:       make(chan error, 1),
		memory: &slaveMemory{
			coils:            make(map[uint16]bool),
			discreteInputs:   make(map[uint16]bool),
			holdingRegisters: make(map[uint16]uint16),
			inputRegisters:   make(map[uint16]uint16),
		},
	}
}

func (s *slaveTCPServer) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		return nil
	}

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	s.listener = listener
	go s.acceptLoop(listener)
	return nil
}

func (s *slaveTCPServer) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var closeErr error
	if s.listener != nil {
		closeErr = s.listener.Close()
		s.listener = nil
	}

	for conn := range s.connections {
		if err := conn.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
		delete(s.connections, conn)
	}

	return closeErr
}

func (s *slaveTCPServer) Errors() <-chan error {
	return s.errCh
}

func (s *slaveTCPServer) PublishTags() {
	tagIDs := make([]uint, 0, len(s.tags))
	for _, tag := range s.tags {
		if tag != nil {
			tagIDs = append(tagIDs, tag.ID)
		}
	}

	values, err := ReadTagsSlave(s, s.tags, nil)
	if err != nil {
		jobs.ApplyPoll(s.clpID, tagIDs, values, time.Now())
		return
	}
	jobs.ApplyPoll(s.clpID, tagIDs, values, time.Now())
}

func (s *slaveTCPServer) ReadCoils(address, quantity uint16) ([]byte, error) {
	s.memory.mu.RLock()
	defer s.memory.mu.RUnlock()

	result := make([]byte, bytesForBits(quantity))
	for i := uint16(0); i < quantity; i++ {
		if s.memory.coils[address+i] {
			result[i/8] |= 1 << (i % 8)
		}
	}
	return result, nil
}

func (s *slaveTCPServer) ReadDiscreteInputs(address, quantity uint16) ([]byte, error) {
	s.memory.mu.RLock()
	defer s.memory.mu.RUnlock()

	result := make([]byte, bytesForBits(quantity))
	for i := uint16(0); i < quantity; i++ {
		if s.memory.discreteInputs[address+i] {
			result[i/8] |= 1 << (i % 8)
		}
	}
	return result, nil
}

func (s *slaveTCPServer) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	return s.readRegisters(s.memory.holdingRegisters, address, quantity)
}

func (s *slaveTCPServer) ReadInputRegisters(address, quantity uint16) ([]byte, error) {
	return s.readRegisters(s.memory.inputRegisters, address, quantity)
}

func (s *slaveTCPServer) WriteSingleCoil(address, value uint16) ([]byte, error) {
	if value != 0x0000 && value != 0xFF00 {
		return nil, fmt.Errorf("valor de coil invalido: %d", value)
	}

	s.memory.mu.Lock()
	s.memory.coils[address] = value == 0xFF00
	s.memory.mu.Unlock()

	s.PublishTags()
	response := make([]byte, 4)
	binary.BigEndian.PutUint16(response[0:2], address)
	binary.BigEndian.PutUint16(response[2:4], value)
	return response, nil
}

func (s *slaveTCPServer) WriteSingleRegister(address, value uint16) ([]byte, error) {
	s.memory.mu.Lock()
	s.memory.holdingRegisters[address] = value
	s.memory.mu.Unlock()

	s.PublishTags()
	response := make([]byte, 4)
	binary.BigEndian.PutUint16(response[0:2], address)
	binary.BigEndian.PutUint16(response[2:4], value)
	return response, nil
}

func (s *slaveTCPServer) WriteMultipleCoils(address, quantity uint16, value []byte) ([]byte, error) {
	if int(bytesForBits(quantity)) > len(value) {
		return nil, fmt.Errorf("dados insuficientes para coils")
	}

	s.memory.mu.Lock()
	for i := uint16(0); i < quantity; i++ {
		s.memory.coils[address+i] = value[i/8]&(1<<(i%8)) != 0
	}
	s.memory.mu.Unlock()

	s.PublishTags()
	return writeMultipleResponse(address, quantity), nil
}

func (s *slaveTCPServer) WriteMultipleRegisters(address, quantity uint16, value []byte) ([]byte, error) {
	if int(quantity)*2 > len(value) {
		return nil, fmt.Errorf("dados insuficientes para registers")
	}

	s.memory.mu.Lock()
	for i := uint16(0); i < quantity; i++ {
		offset := int(i) * 2
		s.memory.holdingRegisters[address+i] = binary.BigEndian.Uint16(value[offset : offset+2])
	}
	s.memory.mu.Unlock()

	s.PublishTags()
	return writeMultipleResponse(address, quantity), nil
}

func (s *slaveTCPServer) readRegisters(registers map[uint16]uint16, address, quantity uint16) ([]byte, error) {
	s.memory.mu.RLock()
	defer s.memory.mu.RUnlock()

	result := make([]byte, int(quantity)*2)
	for i := uint16(0); i < quantity; i++ {
		binary.BigEndian.PutUint16(result[int(i)*2:], registers[address+i])
	}
	return result, nil
}

func (s *slaveTCPServer) acceptLoop(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if s.isCurrentListener(listener) && !errors.Is(err, net.ErrClosed) {
				s.reportError(err)
			}
			return
		}

		s.mu.Lock()
		s.connections[conn] = struct{}{}
		s.mu.Unlock()

		go s.handleConnection(conn)
	}
}

func (s *slaveTCPServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.mu.Lock()
		delete(s.connections, conn)
		s.mu.Unlock()
	}()

	for {
		if s.timeout > 0 {
			if err := conn.SetDeadline(time.Now().Add(s.timeout)); err != nil {
				return
			}
		}

		header := make([]byte, 7)
		if _, err := io.ReadFull(conn, header); err != nil {
			return
		}

		length := int(binary.BigEndian.Uint16(header[4:6]))
		if binary.BigEndian.Uint16(header[2:4]) != 0 || length < 2 {
			return
		}

		body := make([]byte, length-1)
		if _, err := io.ReadFull(conn, body); err != nil {
			return
		}

		unitID := header[6]
		pdu := body
		responsePDU := s.handlePDU(pdu)
		if err := writeModbusTCPResponse(conn, header, unitID, responsePDU); err != nil {
			return
		}
	}
}

func (s *slaveTCPServer) handlePDU(pdu []byte) []byte {
	if len(pdu) == 0 {
		return exceptionResponse(0, modbusExceptionIllegalFunction)
	}

	functionCode := pdu[0]
	switch functionCode {
	case modbusFunctionReadCoils, modbusFunctionReadDiscreteInputs:
		return s.handleReadBits(functionCode, pdu)
	case modbusFunctionReadHoldingRegisters, modbusFunctionReadInputRegisters:
		return s.handleReadRegisters(functionCode, pdu)
	case 5:
		return s.handleWriteSingleCoil(pdu)
	case 6:
		return s.handleWriteSingleRegister(pdu)
	case 15:
		return s.handleWriteMultipleCoils(pdu)
	case 16:
		return s.handleWriteMultipleRegisters(pdu)
	default:
		return exceptionResponse(functionCode, modbusExceptionIllegalFunction)
	}
}

func (s *slaveTCPServer) handleReadBits(functionCode byte, pdu []byte) []byte {
	if len(pdu) < 5 {
		return exceptionResponse(functionCode, modbusExceptionIllegalDataValue)
	}

	address := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])
	if quantity == 0 || int(quantity) > maxCoilsPerRead {
		return exceptionResponse(functionCode, modbusExceptionIllegalDataValue)
	}

	var data []byte
	var err error
	if functionCode == modbusFunctionReadCoils {
		data, err = s.ReadCoils(address, quantity)
	} else {
		data, err = s.ReadDiscreteInputs(address, quantity)
	}
	if err != nil {
		return exceptionResponse(functionCode, modbusExceptionServerFailure)
	}

	return append([]byte{functionCode, byte(len(data))}, data...)
}

func (s *slaveTCPServer) handleReadRegisters(functionCode byte, pdu []byte) []byte {
	if len(pdu) < 5 {
		return exceptionResponse(functionCode, modbusExceptionIllegalDataValue)
	}

	address := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])
	if quantity == 0 || int(quantity) > maxRegistersPerRead {
		return exceptionResponse(functionCode, modbusExceptionIllegalDataValue)
	}

	var data []byte
	var err error
	if functionCode == modbusFunctionReadInputRegisters {
		data, err = s.ReadInputRegisters(address, quantity)
	} else {
		data, err = s.ReadHoldingRegisters(address, quantity)
	}
	if err != nil {
		return exceptionResponse(functionCode, modbusExceptionServerFailure)
	}

	return append([]byte{functionCode, byte(len(data))}, data...)
}

func (s *slaveTCPServer) handleWriteSingleCoil(pdu []byte) []byte {
	if len(pdu) < 5 {
		return exceptionResponse(5, modbusExceptionIllegalDataValue)
	}

	address := binary.BigEndian.Uint16(pdu[1:3])
	value := binary.BigEndian.Uint16(pdu[3:5])
	response, err := s.WriteSingleCoil(address, value)
	if err != nil {
		return exceptionResponse(5, modbusExceptionIllegalDataValue)
	}
	return append([]byte{5}, response...)
}

func (s *slaveTCPServer) handleWriteSingleRegister(pdu []byte) []byte {
	if len(pdu) < 5 {
		return exceptionResponse(6, modbusExceptionIllegalDataValue)
	}

	address := binary.BigEndian.Uint16(pdu[1:3])
	value := binary.BigEndian.Uint16(pdu[3:5])
	response, err := s.WriteSingleRegister(address, value)
	if err != nil {
		return exceptionResponse(6, modbusExceptionServerFailure)
	}
	return append([]byte{6}, response...)
}

func (s *slaveTCPServer) handleWriteMultipleCoils(pdu []byte) []byte {
	if len(pdu) < 6 {
		return exceptionResponse(15, modbusExceptionIllegalDataValue)
	}

	address := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])
	byteCount := int(pdu[5])
	if quantity == 0 || int(quantity) > maxCoilsPerRead || len(pdu) < 6+byteCount || byteCount != int(bytesForBits(quantity)) {
		return exceptionResponse(15, modbusExceptionIllegalDataValue)
	}

	response, err := s.WriteMultipleCoils(address, quantity, pdu[6:6+byteCount])
	if err != nil {
		return exceptionResponse(15, modbusExceptionServerFailure)
	}
	return append([]byte{15}, response...)
}

func (s *slaveTCPServer) handleWriteMultipleRegisters(pdu []byte) []byte {
	if len(pdu) < 6 {
		return exceptionResponse(16, modbusExceptionIllegalDataValue)
	}

	address := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])
	byteCount := int(pdu[5])
	if quantity == 0 || int(quantity) > maxRegistersPerRead || len(pdu) < 6+byteCount || byteCount != int(quantity)*2 {
		return exceptionResponse(16, modbusExceptionIllegalDataValue)
	}

	response, err := s.WriteMultipleRegisters(address, quantity, pdu[6:6+byteCount])
	if err != nil {
		return exceptionResponse(16, modbusExceptionServerFailure)
	}
	return append([]byte{16}, response...)
}

func (s *slaveTCPServer) isCurrentListener(listener net.Listener) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listener == listener
}

func (s *slaveTCPServer) reportError(err error) {
	select {
	case s.errCh <- err:
	default:
	}
}

func bytesForBits(quantity uint16) uint16 {
	return (quantity + 7) / 8
}

func exceptionResponse(functionCode, exceptionCode byte) []byte {
	return []byte{functionCode | 0x80, exceptionCode}
}

func writeMultipleResponse(address, quantity uint16) []byte {
	response := make([]byte, 4)
	binary.BigEndian.PutUint16(response[0:2], address)
	binary.BigEndian.PutUint16(response[2:4], quantity)
	return response
}

func writeModbusTCPResponse(conn net.Conn, requestHeader []byte, unitID byte, pdu []byte) error {
	response := make([]byte, 7+len(pdu))
	copy(response[0:4], requestHeader[0:4])
	binary.BigEndian.PutUint16(response[4:6], uint16(len(pdu)+1))
	response[6] = unitID
	copy(response[7:], pdu)

	_, err := conn.Write(response)
	return err
}
