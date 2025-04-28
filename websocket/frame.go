package websocket

import (
	"bufio"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

const (
	TEXT   = 1
	BINARY = 2
	CLOSE  = 8
	PING   = 9
	PONG   = 10
)

var (
	OPCODE0   = errors.New("OPCODE 0")
	OPCODE1   = errors.New("OPCODE 1")
	OPCODE2   = errors.New("OPCODE 2")
	OPCODE8   = errors.New("OPCODE 8")
	OPCODE9   = errors.New("OPCODE 9")
	OPCODE10  = errors.New("OPCODE 10")
	OPCODE3_7 = errors.New("OPCODE 3-7")
	FINNOT1   = errors.New("FIN NOT 1")
	RSVNOT0   = errors.New("RSV NOT 0")
	CLOSED    = errors.New("CLOSED")
)

// ReadFrame read bytes from reader
func ReadFrame(reader *bufio.Reader) ([]byte, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	fin := firstByte&0x80 == 0x80
	rsv1 := firstByte&0x40 == 0x40
	rsv2 := firstByte&0x20 == 0x20
	rsv3 := firstByte&0x10 == 0x10
	opcode := firstByte & 0x0F
	if !fin {
		return nil, FINNOT1
	}
	switch opcode {
	case 0:
		return nil, OPCODE0
	case 3, 4, 5, 6, 7:
		return nil, OPCODE3_7
	case 8:
		return nil, OPCODE8
	case 9:
		return nil, OPCODE9
	case 10:
		return nil, OPCODE10
	}
	if rsv1 || rsv2 || rsv3 {
		return nil, RSVNOT0
	}

	secondByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	masked := secondByte&0x80 == 0x80
	payloadLength := int(secondByte & 0x7F)

	if payloadLength == 126 {

		lengthBytes := make([]byte, 2)
		_, err := reader.Read(lengthBytes)
		if err != nil {
			return nil, err
		}
		payloadLength = int(binary.BigEndian.Uint16(lengthBytes))
	} else if payloadLength == 127 {

		lengthBytes := make([]byte, 8)
		_, err := reader.Read(lengthBytes)
		if err != nil {
			return nil, err
		}
		payloadLength = int(binary.BigEndian.Uint64(lengthBytes))
	}

	mask := make([]byte, 4)
	if masked {
		_, err := reader.Read(mask)
		if err != nil {
			return nil, err
		}
	}
	payloadData := make([]byte, payloadLength)
	_, err = reader.Read(payloadData)
	if err != nil {
		return nil, err
	}

	if masked {
		for i := 0; i < payloadLength; i++ {
			payloadData[i] ^= mask[i%4]
		}
	}

	switch opcode {
	case 1:
		return payloadData, OPCODE1
	case 2:
		return payloadData, OPCODE2
	}
	return payloadData, nil
}

// WriteFrame write data to conn with datatype
func WriteFrame(conn net.Conn, data []byte, datatype int) error {

	if datatype == CLOSE {
		return conn.Close()
	}

	frameLength := len(data)

	// Write the first byte: FIN flag and OPCODE
	firstByte := byte(0x80)
	opcode := byte(datatype)
	firstByte |= opcode
	err := binary.Write(conn, binary.BigEndian, &firstByte)
	if err != nil {
		return err
	}

	// Write the second byte: payload length
	secondByte := byte(frameLength)
	err = binary.Write(conn, binary.BigEndian, &secondByte)
	if err != nil {
		return err
	}

	// Write the payload data
	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Ping sends ping to connection and waits for pong
func Ping(conn net.Conn, bt []byte, d time.Duration) error {
	err := WriteFrame(conn, bt, TEXT)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)

	er := make(chan error)
	go func() {
		_, err := ReadFrame(reader)
		er <- err
	}()

	select {
	case err := <-er:
		if errors.Is(err, OPCODE10) {
			return nil
		}
		return errors.New("NOT PONG")
	case <-time.After(d):
		return errors.New("PING TIMEOUT")
	}
}

// Pone sends pong to connection
func Pone(conn net.Conn) error {
	data, err := ReadFrame(bufio.NewReader(conn))
	if !errors.Is(err, OPCODE9) {
		return errors.New("NOT PING")
	}
	err = WriteFrame(conn, data, PONG)
	if err != nil {
		return err
	}
	return nil
}
