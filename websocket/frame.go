package websocket

import (
	"encoding/binary"
)

// ReadFrame read bytes from reader, return data, datatype, error.
func ReadFrame(reader Reader) (data []byte, dadatype int, err error) {
	dadatype = -1
	firstByte, err := reader.ReadByte()
	if err != nil {
		return
	}
	fin := firstByte&0x80 == 0x80
	rsv1 := firstByte&0x40 == 0x40
	rsv2 := firstByte&0x20 == 0x20
	rsv3 := firstByte&0x10 == 0x10
	opcode := firstByte & 0x0F
	dadatype = int(opcode)
	if !fin {
		err = ErrFinNot1
		return
	}
	switch opcode {
	case 0:
		err = ErrOpcode0
		return
	case 3, 4, 5, 6, 7:
		err = ErrOpcode37
		return
	case 8, 9, 10:
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

// WriteFrame write data with datatype to writer.
func WriteFrame(writer Writer, data []byte, datatype int) error {
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
	return err
}

// CloseFrame close the connection with status code and reason.
func CloseFrame(writer Writer, code int, reason string) error {
	var payload []byte
	if code != 0 {
		payload = make([]byte, 0, 2+len(reason))
		payload = append(payload, byte(code>>8), byte(code))
		if reason != "" {
			payload = append(payload, []byte(reason)...)
		}
	}
	return WriteFrame(writer, payload, CLOSE)
}
