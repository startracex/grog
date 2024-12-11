package websocket

type WSGroup struct {
	list []*WS
}

// NewWSGroup return empty WSGroup
func NewWSGroup() *WSGroup {
	return &WSGroup{}
}

// Add connection
func (wsg *WSGroup) Add(ws ...*WS) {
	wsg.list = append(wsg.list, ws...)
}

// Close all connections
func (wsg *WSGroup) Close() {
	for _, ws := range wsg.list {
		ws.Close()
	}
	wsg.list = []*WS{}
}

// Send data to all connections
func (wsg *WSGroup) Send(data []byte, datatype int) {
	for _, ws := range wsg.list {
		ws.Send(data, datatype)
	}
}

// Message wait for any message
func (wsg *WSGroup) Message() []byte {
	dataCh := make(chan []byte)
	for i, ws := range wsg.list {
		index := i
		go func(ws *WS) {
			data, err := ws.Message()
			if err != nil {
				ws.Close()
				wsg.Splice(index, 1)
			}
			dataCh <- data

		}(ws)
	}
	return <-dataCh
}

// Clean remove connection which has Closed:true, return removed length
func (wsg *WSGroup) Clean() int {
	initialLen := wsg.Len()
	for i := initialLen - 1; i >= 0; i-- {
		if wsg.list[i].Closed {
			wsg.list = append(wsg.list[:i], wsg.list[i+1:]...)
		}
	}
	return initialLen - wsg.Len()
}

// Splice remove connection from index to index+count, returns removed length
func (wsg *WSGroup) Splice(index, count int) int {
	initialLen := wsg.Len()
	wsg.list = append(wsg.list[:index], wsg.list[index+count:]...)
	return initialLen - wsg.Len()
}

// Len return length of list
func (wsg *WSGroup) Len() int {
	return len(wsg.list)
}
