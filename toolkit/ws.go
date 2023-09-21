package toolkit

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

const A = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

const (
	TEXT = 1

	BINARY = 2

	CLOSE = 8

	PING = 9

	PONG = 10
)

var (
	OPCODE0  = errors.New("OPCODE 0")
	OPCODE1  = errors.New("OPCODE 1")
	OPCODE2  = errors.New("OPCODE 2")
	OPCODE8  = errors.New("OPCODE 8")
	OPCODE9  = errors.New("OPCODE 9")
	OPCODE10 = errors.New("OPCODE 10")

	OPCODE3_7 = errors.New("OPCODE 3-7")

	FINNOT1 = errors.New("FIN NOT 1")

	RSVNOT0 = errors.New("RSV NOT 0")
)

var CLOSED = errors.New("CLOSED")

// IsWebsocketUpgradeRequest check if request has websocket header ("Upgrade": "websocket")
func IsWebsocketUpgradeRequest(r *http.Request) bool {
	return r.Header.Get("Upgrade") == "websocket"
}

func GetKey(r *http.Request) string {
	return r.Header.Get("Sec-WebSocket-Key")
}

// GetAccept get accept from key
func GetAccept(key string) string {
	sha := sha1.New()
	sha.Write([]byte(key + A))
	s1 := sha.Sum(nil)
	return base64.StdEncoding.EncodeToString(s1)
}

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

type WS struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Closed bool
}

// NewWS return empty WS
func NewWS() *WS {
	return &WS{}
}

// Upgrade http connection to websocket
func (ws *WS) Upgrade(w http.ResponseWriter, r *http.Request) {
	if IsWebsocketUpgradeRequest(r) {
		key := GetKey(r)
		Accept := GetAccept(key)
		hijack := w.(http.Hijacker)
		conn, _, err := hijack.Hijack()
		if err != nil {
			return
		}
		_, err = conn.Write([]byte("HTTP/1.1 101 Switching Protocols\nUpgrade: websocket\nConnection: Upgrade\nSec-WebSocket-Accept: " + Accept + "\n\n"))
		if err != nil {
			return
		}
		ws.Conn = conn
		ws.Reader = bufio.NewReader(conn)
	}
}

// Upgrade converse to websocket instance
func Upgrade(w http.ResponseWriter, r *http.Request) *WS {
	ws := NewWS()
	ws.Upgrade(w, r)
	return ws
}

// Send data to conn with datatype
func (ws *WS) Send(data []byte, datatype int) error {
	if ws.Closed {
		return CLOSED
	}
	return WriteFrame(ws.Conn, data, datatype)
}

// Message wait for a message
func (ws *WS) Message() ([]byte, error) {
	data, err := ReadFrame(ws.Reader)
	if errors.Is(err, OPCODE8) || errors.Is(err, io.EOF) { // OPCODE8 / EOF close
		ws.Closed = true
		return nil, CLOSED
	} else if errors.Is(err, OPCODE9) { // OPCODE9 write PONG
		err := WriteFrame(ws.Conn, data, PONG)
		if err != nil {
			return nil, err
		}
		return data, OPCODE9
	} else { // OPCODE Other
		return data, nil
	}
}

// Close connection
func (ws *WS) Close() error {
	return ws.Conn.Close()
}

func (ws *WS) Ping(bt []byte, d time.Duration) error {
	if ws.Closed {
		return CLOSED
	}
	return Ping(ws.Conn, bt, d)
}

func (ws *WS) Pone() error {
	return Pone(ws.Conn)
}

// IsOpen return ws.Closed equals true
func (ws *WS) IsOpen() bool {
	return !ws.Closed
}

type WSGroup struct {
	list []*WS
}

// NewWSGroup return empty WSGroup
func NewWSGroup() *WSGroup {
	return &WSGroup{}
}

// Add connection to group
func (w *WSGroup) Add(ws *WS) {
	w.list = append(w.list, ws)
}

// Send data to all connections
func (w *WSGroup) Send(data []byte, t int) error {
	var err error = nil
	for _, ws := range w.list {
		err = ws.Send(data, t)
		if err != nil {
			ws.Closed = true
		}
	}
	return err
}

// Message wait for any message
func (w *WSGroup) Message() ([]byte, error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	for _, ws := range w.list {
		go func(ws *WS) {
			data, err := ws.Message()
			if err != nil {
				ws.Closed = true
			}
			dataCh <- data
			errCh <- err
		}(ws)
	}
	return <-dataCh, <-errCh
}

// Clean remove connection which has Closed:true
func (w *WSGroup) Clean() int {
	f := w.Len()
	for i := f - 1; i >= 0; i-- {
		if w.list[i].Closed {
			w.list = append(w.list[:i], w.list[i+1:]...)
		}
	}
	return f - w.Len()
}

func (w *WSGroup) Len() int {
	return len(w.list)
}
