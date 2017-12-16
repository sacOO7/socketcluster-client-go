package listener

type Listener struct {
	emitAckListener map[int]func(eventName string, error interface{}, data interface{})
	onListener      map[int]func(eventName string, data interface{})
	onAckListener map[int]func(eventName string, data interface{}, ack func(error interface{}, data interface{}))
}

func Init() Listener {
	return Listener {
		emitAckListener: make(map[int]func(eventName string, error interface{}, data interface{})),
		onListener: make(map[int]func(eventName string, data interface{})),
		onAckListener: make(map[int]func(eventName string, data interface{}, ack func(error interface{}, data interface{}))),
	}
}

func (listener *Listener) putEmitAck(id int, ack func(eventName string, error interface{}, data interface{})) {
	listener.emitAckListener[id] = ack
}

func (listener *Listener) executeEmitAck(id int, eventName string, error interface{}, data interface{}) {
	ack := listener.emitAckListener[id];
	if ack != nil {
		ack(eventName, error, data);
	}
}

func (listener *Listener) putOnListener(id int, onListener func(eventName string, data interface{})) {
	listener.onListener[id] = onListener
}

func (listener *Listener) executeOnListener(id int, eventName string, data interface{}) {
	on := listener.onListener[id];
	if on != nil {
		on(eventName, data);
	}
}

func (listener *Listener) putOnAckListener(id int, onAckListener func(eventName string, data interface{}, ack func(error interface{}, data interface{}))) {
	listener.onAckListener[id] = onAckListener
}

func (listener *Listener) executeOnAckListener(id int, eventName string, data interface{}, ack func(error interface{}, data interface{})) {
	onAck := listener.onAckListener[id]
	if onAck != nil {
		onAck(eventName, data, ack)
	}
}
