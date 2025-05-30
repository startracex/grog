package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

var ErrClosed = errors.New("CLOSED")

func GetKey(h http.Header) string {
	return h.Get("Sec-WebSocket-Key")
}

func SumAccept(key string) string {
	sha := sha1.New()
	sha.Write([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(sha.Sum(nil))
}

func IsUpgradeWebsocket(h http.Header) bool {
	return h.Get("Upgrade") == "websocket"
}

func WriteAccept(w io.Writer, accept string) (int, error) {
	return w.Write([]byte("HTTP/1.1 101 Switching Protocols\nUpgrade: websocket\nConnection: Upgrade\nSec-WebSocket-Accept: " + accept + "\n\n"))
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
	if IsUpgradeWebsocket(r.Header) {
		key := GetKey(r.Header)
		accept := SumAccept(key)
		hijack := w.(http.Hijacker)
		conn, _, err := hijack.Hijack()
		if err != nil {
			return
		}
		_, err = WriteAccept(conn, accept)
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
		return ErrClosed
	}
	return WriteFrame(ws.Conn, data, datatype)
}

// Message wait for a message
func (ws *WS) Message() ([]byte, error) {
	data, err := ReadFrame(ws.Reader)
	if IsClose(err) { // OPCODE 8 close / EOF close
		ws.Closed = true
		return nil, err
	}
	if IsPing(err) { // OPCODE 9 ping
		err := WriteFrame(ws.Conn, data, PONG)
		if err != nil {
			return nil, err
		}
		return data, ErrOpcode9
	}

	return data, nil
}

// Close connection
func (ws *WS) Close() error {
	if ws.Closed {
		return ErrClosed
	} else {
		ws.Closed = true
		return ws.Conn.Close()
	}
}

func (ws *WS) Ping(bt []byte, d time.Duration) error {
	if ws.Closed {
		return ErrClosed
	}
	return Ping(ws.Conn, bt, d)
}

func (ws *WS) Pone(bt []byte) error {
	return Pone(ws.Conn, bt)
}
