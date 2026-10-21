package modbusSlave

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type modbusSlaveClient interface {
	ReadCoils(address, quantity uint16) ([]byte, error)
	ReadDiscreteInputs(address, quantity uint16) ([]byte, error)
	ReadHoldingRegisters(address, quantity uint16) ([]byte, error)
	ReadInputRegisters(address, quantity uint16) ([]byte, error)
	WriteSingleCoil(address, value uint16) ([]byte, error)
	WriteSingleRegister(address, value uint16) ([]byte, error)
	WriteMultipleCoils(address, quantity uint16, value []byte) ([]byte, error)
	WriteMultipleRegisters(address, quantity uint16, value []byte) ([]byte, error)
}

type slaveConnector interface {
	Connect() error
	Close() error
}

type slaveTCPClient struct {
	address     string
	slaveID     byte
	timeout     time.Duration
	idleTimeout time.Duration

	mu       sync.Mutex
	conn     net.Conn
	txID     uint16
	lastUsed time.Time
}

type modbusExceptionError struct {
	functionCode  byte
	exceptionCode byte
}

func (e modbusExceptionError) Error() string {
	return fmt.Sprintf("modbus exception function=%d code=%d", e.functionCode, e.exceptionCode)
}

func newSlaveTCPClient(address string, slaveID byte, timeout, idleTimeout time.Duration) *slaveTCPClient {
	return &slaveTCPClient{
		address:     address,
		slaveID:     slaveID,
		timeout:     timeout,
		idleTimeout: idleTimeout,
	}
}

func (c *slaveTCPClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.connectLocked()
}

func (c *slaveTCPClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.closeLocked()
}

func (c *slaveTCPClient) ReadCoils(address, quantity uint16) ([]byte, error) {
	return c.readBits(1, address, quantity)
}

func (c *slaveTCPClient) ReadDiscreteInputs(address, quantity uint16) ([]byte, error) {
	return c.readBits(2, address, quantity)
}

func (c *slaveTCPClient) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	return c.readRegisters(3, address, quantity)
}

func (c *slaveTCPClient) ReadInputRegisters(address, quantity uint16) ([]byte, error) {
	return c.readRegisters(4, address, quantity)
}

func (c *slaveTCPClient) WriteSingleCoil(address, value uint16) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint16(payload[0:2], address)
	binary.BigEndian.PutUint16(payload[2:4], value)
	return c.transact(5, payload)
}

func (c *slaveTCPClient) WriteSingleRegister(address, value uint16) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint16(payload[0:2], address)
	binary.BigEndian.PutUint16(payload[2:4], value)
	return c.transact(6, payload)
}

func (c *slaveTCPClient) WriteMultipleCoils(address, quantity uint16, value []byte) ([]byte, error) {
	payload := make([]byte, 5+len(value))
	binary.BigEndian.PutUint16(payload[0:2], address)
	binary.BigEndian.PutUint16(payload[2:4], quantity)
	payload[4] = byte(len(value))
	copy(payload[5:], value)
	return c.transact(15, payload)
}

func (c *slaveTCPClient) WriteMultipleRegisters(address, quantity uint16, value []byte) ([]byte, error) {
	payload := make([]byte, 5+len(value))
	binary.BigEndian.PutUint16(payload[0:2], address)
	binary.BigEndian.PutUint16(payload[2:4], quantity)
	payload[4] = byte(len(value))
	copy(payload[5:], value)
	return c.transact(16, payload)
}

func (c *slaveTCPClient) readBits(functionCode byte, address, quantity uint16) ([]byte, error) {
	response, err := c.read(functionCode, address, quantity)
	if err != nil {
		return nil, err
	}

	expected := int(response[0])
	if len(response) < 1+expected {
		return nil, fmt.Errorf("resposta modbus invalida para funcao %d", functionCode)
	}

	return append([]byte(nil), response[1:1+expected]...), nil
}

func (c *slaveTCPClient) readRegisters(functionCode byte, address, quantity uint16) ([]byte, error) {
	response, err := c.read(functionCode, address, quantity)
	if err != nil {
		return nil, err
	}

	expected := int(response[0])
	if len(response) < 1+expected {
		return nil, fmt.Errorf("resposta modbus invalida para funcao %d", functionCode)
	}

	return append([]byte(nil), response[1:1+expected]...), nil
}

func (c *slaveTCPClient) read(functionCode byte, address, quantity uint16) ([]byte, error) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint16(payload[0:2], address)
	binary.BigEndian.PutUint16(payload[2:4], quantity)

	response, err := c.transact(functionCode, payload)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 {
		return nil, fmt.Errorf("resposta modbus vazia para funcao %d", functionCode)
	}

	return response, nil
}

func (c *slaveTCPClient) transact(functionCode byte, payload []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.ensureConnectedLocked(); err != nil {
		return nil, err
	}

	c.txID++
	txID := c.txID

	pdu := append([]byte{functionCode}, payload...)
	request := make([]byte, 7+len(pdu))
	binary.BigEndian.PutUint16(request[0:2], txID)
	binary.BigEndian.PutUint16(request[2:4], 0)
	binary.BigEndian.PutUint16(request[4:6], uint16(len(pdu)+1))
	request[6] = c.slaveID
	copy(request[7:], pdu)

	if c.timeout > 0 {
		if err := c.conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
			_ = c.closeLocked()
			return nil, err
		}
	}

	if _, err := c.conn.Write(request); err != nil {
		_ = c.closeLocked()
		return nil, err
	}

	header := make([]byte, 7)
	if _, err := io.ReadFull(c.conn, header); err != nil {
		_ = c.closeLocked()
		return nil, err
	}

	responseTxID := binary.BigEndian.Uint16(header[0:2])
	if responseTxID != txID {
		_ = c.closeLocked()
		return nil, fmt.Errorf("transaction id inesperado: recebido=%d esperado=%d", responseTxID, txID)
	}

	if binary.BigEndian.Uint16(header[2:4]) != 0 {
		_ = c.closeLocked()
		return nil, fmt.Errorf("resposta modbus com protocolo invalido")
	}

	length := int(binary.BigEndian.Uint16(header[4:6]))
	if length < 2 {
		_ = c.closeLocked()
		return nil, fmt.Errorf("resposta modbus com tamanho invalido")
	}

	body := make([]byte, length)
	if _, err := io.ReadFull(c.conn, body); err != nil {
		_ = c.closeLocked()
		return nil, err
	}

	if body[0] != c.slaveID {
		_ = c.closeLocked()
		return nil, fmt.Errorf("slave id inesperado: recebido=%d esperado=%d", body[0], c.slaveID)
	}

	response := body[1:]
	if len(response) == 0 {
		_ = c.closeLocked()
		return nil, fmt.Errorf("pdu modbus vazia")
	}

	if response[0] == functionCode|0x80 {
		if len(response) < 2 {
			_ = c.closeLocked()
			return nil, fmt.Errorf("resposta de excecao modbus invalida")
		}
		return nil, modbusExceptionError{
			functionCode:  functionCode,
			exceptionCode: response[1],
		}
	}

	if response[0] != functionCode {
		_ = c.closeLocked()
		return nil, fmt.Errorf("codigo de funcao inesperado: recebido=%d esperado=%d", response[0], functionCode)
	}

	c.lastUsed = time.Now()
	return append([]byte(nil), response[1:]...), nil
}

func (c *slaveTCPClient) ensureConnectedLocked() error {
	if c.conn == nil {
		return c.connectLocked()
	}

	if c.idleTimeout > 0 && !c.lastUsed.IsZero() && time.Since(c.lastUsed) > c.idleTimeout {
		if err := c.closeLocked(); err != nil {
			return err
		}
		return c.connectLocked()
	}

	return nil
}

func (c *slaveTCPClient) connectLocked() error {
	if err := c.closeLocked(); err != nil {
		return err
	}

	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.conn = conn
	c.lastUsed = time.Now()
	return nil
}

func (c *slaveTCPClient) closeLocked() error {
	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.conn = nil
	c.lastUsed = time.Time{}
	return err
}
