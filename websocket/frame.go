package websocket

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
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
	ErrOpcode0  = errors.New("OPCODE 0")
	ErrOpcode1  = errors.New("OPCODE 1")
	ErrOpcode2  = errors.New("OPCODE 2")
	ErrOpcode37 = errors.New("OPCODE 3-7")
	ErrOpcode8  = errors.New("OPCODE 8")
	ErrOpcode9  = errors.New("OPCODE 9")
	ErrOpcodeA  = errors.New("OPCODE 10")
	ErrOpcodeBF = errors.New("OPCODE 11-15")
	ErrFinNot1  = errors.New("FIN NOT 1")
	ErrRsvNot0  = errors.New("RSV NOT 0")
)

// ReadFrame read bytes from reader, if datatype is 1 or 2, error is nil.
func ReadFrame(reader *bufio.Reader) ([]byte, error) {
	data, code, err := ReadTypeFrame(reader)
	if code == 1 || code == 2 {
		return data, nil
	}
	return nil, err
}

// ReadTypeFrame read bytes from reader, return data, datatype, error.
func ReadTypeFrame(reader *bufio.Reader) (data []byte, code int, err error) {
	code = -1
	firstByte, err := reader.ReadByte()
	if err != nil {
		return
	}
	fin := firstByte&0x80 == 0x80
	rsv1 := firstByte&0x40 == 0x40
	rsv2 := firstByte&0x20 == 0x20
	rsv3 := firstByte&0x10 == 0x10
	opcode := firstByte & 0x0F
	code = int(opcode)
	if !fin {
		err = ErrFinNot1
		return
	}
	switch opcode {
	case 0:
		err = ErrOpcode0
		return
	case 1:
		err = ErrOpcode1
	case 2:
		err = ErrOpcode2
	case 3, 4, 5, 6, 7:
		err = ErrOpcode37
		return
	case 8:
		err = ErrOpcode8
		return
	case 9:
		err = ErrOpcode9
		return
	case 10:
		err = ErrOpcodeA
		return
	case 11, 12, 13, 14, 15:
		err = ErrOpcodeBF
		return
	}
	if rsv1 || rsv2 || rsv3 {
		err = ErrRsvNot0
		return
	}

	secondByte, err := reader.ReadByte()
	if err != nil {
		return
	}
	masked := secondByte&0x80 == 0x80
	payloadLength := int(secondByte & 0x7F)

	switch payloadLength {
	case 126:
		lengthBytes := make([]byte, 2)
		_, err = reader.Read(lengthBytes)
		if err != nil {
			return
		}
		payloadLength = int(binary.BigEndian.Uint16(lengthBytes))
	case 127:
		lengthBytes := make([]byte, 8)
		_, err = reader.Read(lengthBytes)
		if err != nil {
			return
		}
		payloadLength = int(binary.BigEndian.Uint64(lengthBytes))
	}

	mask := make([]byte, 4)
	if masked {
		_, err = reader.Read(mask)
		if err != nil {
			return
		}
	}
	payloadData := make([]byte, payloadLength)
	_, err = reader.Read(payloadData)
	if err != nil {
		return
	}

	if masked {
		for i := range payloadLength {
			payloadData[i] ^= mask[i%4]
		}
	}

	data = payloadData
	return
}

func WriteFrame(writer *bufio.Writer, data []byte, datatype int) error {
	if datatype == CLOSE {
		return writer.Flush()
	}

	frameLength := len(data)

	// Write the first byte: FIN flag and OPCODE
	firstByte := byte(0x80)
	opcode := byte(datatype)
	firstByte |= opcode
	err := writer.WriteByte(firstByte)
	if err != nil {
		return err
	}

	// Write the second byte: payload length
	var secondByte byte
	if frameLength < 126 {
		secondByte = byte(frameLength)
		err = writer.WriteByte(secondByte)
	} else if frameLength <= 0xFFFF {
		secondByte = 126
		err = writer.WriteByte(secondByte)
		if err == nil {
			err = binary.Write(writer, binary.BigEndian, uint16(frameLength))
		}
	} else {
		secondByte = 127
		err = writer.WriteByte(secondByte)
		if err == nil {
			err = binary.Write(writer, binary.BigEndian, uint64(frameLength))
		}
	}
	if err != nil {
		return err
	}

	// Write the payload data
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return writer.Flush()
}

// Ping sends ping to connection and waits for pong
func Ping(conn net.Conn, bt []byte, d time.Duration) error {
	writer := bufio.NewWriter(conn)
	err := WriteFrame(writer, bt, PING)
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
		if errors.Is(err, ErrOpcodeA) {
			return nil
		}
		return errors.New("NOT PONG")
	case <-time.After(d):
		return errors.New("PING TIMEOUT")
	}
}

// Pone sends pong to connection
func Pone(conn net.Conn, bt []byte) error {
	writer := bufio.NewWriter(conn)
	err := WriteFrame(writer, bt, PONG)
	if err != nil {
		return err
	}
	return nil
}

func Close(conn net.Conn, bt []byte) error {
	writer := bufio.NewWriter(conn)
	err := WriteFrame(writer, bt, CLOSE)
	if err != nil {
		return err
	}
	return conn.Close()
}

func IsPing(err error) bool {
	return errors.Is(err, ErrOpcode9)
}

func IsPong(err error) bool {
	return errors.Is(err, ErrOpcodeA)
}

func IsClose(err error) bool {
	return errors.Is(err, ErrOpcode8) || errors.Is(err, io.EOF)
}
