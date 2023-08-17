package toolkit

import (
    "bufio"
    "crypto/sha1"
    "encoding/base64"
    "encoding/binary"
    "errors"
    "net"
    "net/http"
)

const A = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

const TEXT = 1
const BINARY = 2

const CLOSE = 8

// IsWebsocketRequest check if request has websocket header ("Upgrade": "websocket")
func IsWebsocketRequest(r *http.Request) bool {
    return r.Header.Get("Upgrade") == "websocket"
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
        return nil, errors.New("FIN NOT 1")
    }
    if opcode == 8 {
        return nil, errors.New("OPCODE 8")
    }
    if rsv1 || rsv2 || rsv3 {
        return nil, errors.New("RSV1/RSV2/RSV3 NOT 0")
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

    return payloadData, nil
}

// WriteFrame write data to conn with datatype
func WriteFrame(conn net.Conn, data []byte, datatype int) error {

    if datatype == CLOSE {
        return conn.Close()
    }

    frameLength := len(data)

    // Write the first byte: FIN flag and opcode
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

type WS struct {
    Conn   net.Conn
    Reader *bufio.Reader
}

func NewWS() *WS {
    return &WS{}
}

// Upgrade http connection to websocket
func (ws *WS) Upgrade(w http.ResponseWriter, r *http.Request) {
    if IsWebsocketRequest(r) {
        key := r.Header.Get("Sec-WebSocket-Key")
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
    return WriteFrame(ws.Conn, data, datatype)
}

// Message wait for message
func (ws *WS) Message() ([]byte, error) {
    ch := make(chan []byte)
    go func() {
        data, err := ReadFrame(ws.Reader)
        if err != nil {
            return
        }
        ch <- data
    }()
    return <-ch, nil
}

// Close connection
func (ws *WS) Close() error {
    return ws.Conn.Close()
}
