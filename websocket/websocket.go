package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
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
	Conn       net.Conn
	Writer     *bufio.Writer
	WriterSize int
	Reader     *bufio.Reader
	ReaderSize int
	Closed     bool
	mu         sync.Mutex
}

// New return empty WS
func New() *WS {
	return &WS{}
}

func (ws *WS) checkSize() {
	if ws.ReaderSize <= 0 {
		ws.ReaderSize = 4096
	}
	if ws.WriterSize <= 0 {
		ws.WriterSize = 4096
	}
}

func (ws *WS) Connect(conn net.Conn) {
	ws.checkSize()
	ws.Conn = conn
	ws.Reader = bufio.NewReaderSize(conn, ws.ReaderSize)
	ws.Writer = bufio.NewWriterSize(conn, ws.WriterSize)
}

var ErrNotUpgrade = errors.New("grog/websocket: not upgrade")

// Upgrade http connection to websocket
func (ws *WS) Upgrade(w http.ResponseWriter, r *http.Request) error {
	if !IsUpgradeWebsocket(r.Header) {
		return ErrNotUpgrade
	}
	key := GetKey(r.Header)
	accept := SumAccept(key)
	hijack := w.(http.Hijacker)
	conn, _, err := hijack.Hijack()
	if err != nil {
		return err
	}
	_, err = WriteAccept(conn, accept)
	if err != nil {
		return err
	}
	ws.Connect(conn)
	return nil
}

// Upgrade converse to websocket instance
func Upgrade(w http.ResponseWriter, r *http.Request) *WS {
	ws := New()
	ws.Upgrade(w, r)
	return ws
}

// Send data to conn with datatype using bufio.Writer
func (ws *WS) Send(data []byte, datatype int) error {
	if ws.Closed {
		return ErrClosed
	}
	return WriteFrame(ws.Writer, data, datatype)
}

// Message waits for a message and handles PING automatically
func (ws *WS) Message() ([]byte, error) {
	data, err := ReadFrame(ws.Reader)
	if IsClose(err) { // OPCODE 8 close / EOF close
		ws.Closed = true
		return nil, err
	}
	if IsPing(err) { // OPCODE 9 ping
		err := WriteFrame(ws.Writer, data, PONG)
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
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.Closed = true
	err := WriteFrame(ws.Writer, nil, CLOSE)
	if err != nil {
		return err
	}

	ws.Closed = true
	return ws.Conn.Close()

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
