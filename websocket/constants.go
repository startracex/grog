package websocket

import "errors"

const (
	TEXT   = 1
	BINARY = 2
	CLOSE  = 8
	PING   = 9
	PONG   = 10
)

const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

var (
	ErrOpcode0     = errors.New("grog/websocket: OPCODE 0")
	ErrOpcode37    = errors.New("grog/websocket: OPCODE 3-7")
	ErrOpcodeBF    = errors.New("grog/websocket: OPCODE 11-15")
	ErrFinNot1     = errors.New("grog/websocket: FIN not 1")
	ErrRsvNot0     = errors.New("grog/websocket: RSV not 0")
	ErrPingTimeout = errors.New("grog/websocket: PING timeout")
	ErrNotUpgrade  = errors.New("grog/websocket: not upgrade")
	ErrClosed      = errors.New("grog/websocket: closed")
)
