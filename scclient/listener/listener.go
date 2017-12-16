package listener

type Listener struct {
	emitAckListener map[int][] interface{}
	onListener      map[int]func(eventName string, data interface{})
	onAckListener map[int]func(eventName string, data interface{}, ack func(error interface{}, data interface{}))
}

func Init() Listener {
	return Listener{
		emitAckListener: make(map[int][] interface{}),
		onListener:      make(map[int]func(eventName string, data interface{})),
		onAckListener: make(map[int]func(eventName string, data interface{}, ack func(error interface{}, data interface{}))),
	}
}

func (listener *Listener) PutEmitAck(id int, eventName string, ack func(eventName string, error interface{}, data interface{})) {
	listener.emitAckListener[id] = [] interface{}{eventName, ack}
}

func (listener *Listener) HandleEmitAck(id int, error interface{}, data interface{}) {
	ackObject := listener.emitAckListener[id];
	if ackObject != nil {
		eventName := ackObject[0].(string)
		ack := ackObject[1].(func(eventName string, error interface{}, data interface{}))
		ack(eventName, error, data);
	}
}

func (listener *Listener) PutOnListener(id int, onListener func(eventName string, data interface{})) {
	listener.onListener[id] = onListener
}

func (listener *Listener) HandleOnListener(id int, eventName string, data interface{}) {
	on := listener.onListener[id];
	if on != nil {
		on(eventName, data);
	}
}

func (listener *Listener) PutOnAckListener(id int, onAckListener func(eventName string, data interface{}, ack func(error interface{}, data interface{}))) {
	listener.onAckListener[id] = onAckListener
}

func (listener *Listener) HandleOnAckListener(id int, eventName string, data interface{}, ack func(error interface{}, data interface{})) {
	onAck := listener.onAckListener[id]
	if onAck != nil {
		onAck(eventName, data, ack)
	}
}
