package listener

type Listener struct {
	emitAckListener map[int][] interface{}
	onListener      map[string]func(eventName string, data interface{})
	onAckListener map[string]func(eventName string, data interface{}, ack func(error interface{}, data interface{}))
}

func Init() Listener {
	return Listener{
		emitAckListener: make(map[int][] interface{}),
		onListener:      make(map[string]func(eventName string, data interface{})),
		onAckListener: make(map[string]func(eventName string, data interface{}, ack func(error interface{}, data interface{}))),
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

func (listener *Listener) PutOnListener(eventName string, onListener func(eventName string, data interface{})) {
	listener.onListener[eventName] = onListener
}

func (listener *Listener) HandleOnListener(eventName string, data interface{}) {
	on := listener.onListener[eventName];
	if on != nil {
		on(eventName, data);
	}
}

func (listener *Listener) PutOnAckListener(eventName string, onAckListener func(eventName string, data interface{}, ack func(error interface{}, data interface{}))) {
	listener.onAckListener[eventName] = onAckListener
}

func (listener *Listener) HandleOnAckListener(eventName string, data interface{}, ack func(error interface{}, data interface{})) {
	onAck := listener.onAckListener[eventName]
	if onAck != nil {
		onAck(eventName, data, ack)
	}
}
