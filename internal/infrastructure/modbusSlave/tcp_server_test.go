package modbusSlave

import (
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"
)

func TestSlaveTCPServerWriteAndReadHoldingRegister(t *testing.T) {
	server := newSlaveTCPServer("127.0.0.1:0", 1, 1, nil, time.Second)
	if err := server.Connect(); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer server.Close()

	address := server.listener.Addr().String()
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		t.Fatalf("DialTimeout() error = %v", err)
	}
	defer conn.Close()

	writePDU := []byte{6, 0, 10, 0, 42}
	writeResponse := sendModbusTCPRequest(t, conn, 1, writePDU)
	if len(writeResponse) != 5 || writeResponse[0] != 6 {
		t.Fatalf("write response = %v, want echo for function 6", writeResponse)
	}

	readPDU := []byte{3, 0, 10, 0, 1}
	readResponse := sendModbusTCPRequest(t, conn, 2, readPDU)
	if len(readResponse) != 4 || readResponse[0] != 3 || readResponse[1] != 2 {
		t.Fatalf("read response = %v, want function 3 with 2 bytes", readResponse)
	}

	if value := binary.BigEndian.Uint16(readResponse[2:4]); value != 42 {
		t.Fatalf("holding register value = %d, want 42", value)
	}
}

func sendModbusTCPRequest(t *testing.T, conn net.Conn, transactionID uint16, pdu []byte) []byte {
	t.Helper()

	request := make([]byte, 7+len(pdu))
	binary.BigEndian.PutUint16(request[0:2], transactionID)
	binary.BigEndian.PutUint16(request[4:6], uint16(len(pdu)+1))
	request[6] = 1
	copy(request[7:], pdu)

	if _, err := conn.Write(request); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	header := make([]byte, 7)
	if _, err := io.ReadFull(conn, header); err != nil {
		t.Fatalf("ReadFull(header) error = %v", err)
	}

	if got := binary.BigEndian.Uint16(header[0:2]); got != transactionID {
		t.Fatalf("transaction id = %d, want %d", got, transactionID)
	}

	bodyLength := int(binary.BigEndian.Uint16(header[4:6])) - 1
	body := make([]byte, bodyLength)
	if _, err := io.ReadFull(conn, body); err != nil {
		t.Fatalf("ReadFull(body) error = %v", err)
	}

	if header[6] != 1 {
		t.Fatalf("unit id = %d, want 1", header[6])
	}
	return body
}
