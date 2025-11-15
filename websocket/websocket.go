package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"net"
	"net/http"
	"sync"
)

type WebSocket struct {
	Conn       net.Conn
	Writer     *bufio.Writer
	WriterSize int
	Reader     *bufio.Reader
	ReaderSize int
	Closed     bool
	mu         sync.Mutex
}

func New() *WebSocket {
	return &WebSocket{}
}

func (ws *WebSocket) Connect(conn net.Conn) {
	if ws.ReaderSize <= 0 {
		ws.ReaderSize = 4096
	}
	if ws.WriterSize <= 0 {
		ws.WriterSize = 4096
	}
	ws.Conn = conn
	ws.Reader = bufio.NewReaderSize(conn, ws.ReaderSize)
	ws.Writer = bufio.NewWriterSize(conn, ws.WriterSize)
}

// Upgrade http connection to websocket
func (ws *WebSocket) Upgrade(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Upgrade") != "websocket" {
		return ErrNotUpgrade
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	sha := sha1.New()
	sha.Write([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	accept := base64.StdEncoding.EncodeToString(sha.Sum(nil))
	hijack := w.(http.Hijacker)
	conn, _, err := hijack.Hijack()
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: " + accept + "\r\n\r\n"))
	if err != nil {
		return err
	}
	ws.Connect(conn)
	return nil
}

// Send data to conn with datatype using bufio.Writer
func (ws *WebSocket) Send(data []byte, datatype int) error {
	if ws.Closed {
		return ErrClosed
	}
	err := WriteFrame(ws.Writer, data, datatype)
	if err != nil {
		return err
	}
	return ws.Writer.Flush()
}

// Message waits for a message
func (ws *WebSocket) Message() ([]byte, int, error) {
	if ws.Closed {
		return nil, -1, ErrClosed
	}
	data, datatype, err := ReadFrame(ws.Reader)
	switch datatype {
	case PING:
		err = ws.Send(nil, PONG)
	case CLOSE:
		err = ws.Close(1000, "")
	}
	return data, datatype, err
}

// Close connection
func (ws *WebSocket) Close(code int, reason string) error {
	if ws.Closed {
		return nil
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.Closed = true
	err := CloseFrame(ws.Writer, code, reason)
	if err != nil {
		return err
	}
	err = ws.Writer.Flush()
	if err != nil {
		return err
	}
	return ws.Conn.Close()
}
