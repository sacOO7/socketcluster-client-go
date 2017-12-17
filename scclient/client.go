package scclient

import (
	"log"
	_ "golang.org/x/net/websocket"
	"github.com/rgamba/evtwebsocket"
	_ "time"
	"github.com/sacOO7/socketcluster-client-go/scclient/models"
	"github.com/sacOO7/socketcluster-client-go/scclient/utils"
	"github.com/sacOO7/socketcluster-client-go/scclient/parser"
)

type Client struct {
	authToken           *string
	url                 string
	counter             utils.AtomicCounter
	socket              *evtwebsocket.Conn
	onConnect           func(client Client)
	onConnectError      func(client Client, err error)
	onDisconnect        func(client Client, err error)
	onSetAuthentication func(client Client, token string)
	onAuthentication    func(client Client, isAuthenticated bool)
	Listener
}

func New(url string) Client {
	return Client{url: url, counter: utils.AtomicCounter{Counter: 0}, Listener: Init()}
}

func (client *Client) SetBasicListener(onConnect func(client Client), onConnectError func(client Client, err error), onDisconnect func(client Client, err error)) {
	client.onConnect = onConnect
	client.onConnectError = onConnectError
	client.onDisconnect = onDisconnect
}

func (client *Client) SetAuthenticationListener(onSetAuthentication func(client Client, token string), onAuthentication func(client Client, isAuthenticated bool)) {
	client.onSetAuthentication = onSetAuthentication
	client.onAuthentication = onAuthentication
}

func (client *Client) registerCallbacks() {
	client.socket = &evtwebsocket.Conn{
		// Fires when the connection is established
		OnConnected: func(w *evtwebsocket.Conn) {
			if client.onConnect != nil {
				client.onConnect(*client)
			}
			client.sendHandshake()
		},
		// Fires when a new message arrives from the server
		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
			log.Printf("%s", msg)

			if utils.IsEqual("#1", msg) {
				w.Send(utils.CreateMessageFromString("#2"));
			} else {
				var messageObject = utils.DeserializeData(msg)
				data, rid, cid, eventname, error := parser.GetMessageDetails(messageObject)

				parseresult := parser.Parse(rid, cid, eventname)

				switch parseresult {
				case parser.ISAUTHENTICATED:
					isAuthenticated := utils.GetIsAuthenticated(messageObject)
					if client.onAuthentication != nil {
						client.onAuthentication(*client, isAuthenticated);
					}
				case parser.SETTOKEN:
					token := utils.GetAuthToken(messageObject)
					if client.onSetAuthentication != nil {
						client.onSetAuthentication(*client, token)
					}

				case parser.REMOVETOKEN:
					client.authToken = nil
				case parser.EVENT:
					if client.hasEventAck(eventname.(string)) {
						client.handleOnAckListener(eventname.(string), data, client.ack(cid))
					} else {
						client.handleOnListener(eventname.(string), data)
					}
				case parser.ACKRECEIVE:
					client.handleEmitAck(rid, error, data)
				case parser.PUBLISH:
					channel := models.GetChannelObject(data)
					client.handleOnListener(channel.Channel, channel.Data)
				}
			}

		},
		// Fires when an error occurs and connection is closed
		OnError: func(err error) {
			if client.onDisconnect != nil {
				client.onDisconnect(*client, err)
			}
		},
	}

}

func (client *Client) Connect() {
	client.registerCallbacks()
	// Connect
	err := client.socket.Dial(client.url, "")
	if err != nil {
		if client.onConnectError != nil {
			client.onConnectError(*client, err)
		}
	}
}

func (client *Client) sendHandshake() {
	handshake := utils.SerializeData(models.GetHandshakeObject(client.authToken, int(client.counter.IncrementAndGet())))
	client.socket.Send(utils.CreateMessageFromByte(handshake));
}

func (client *Client) ack(cid int) func(error interface{}, data interface{}) {
	return func(error interface{}, data interface{}) {
		ackObject := models.GetReceiveEventObject(data, error, cid);
		ackData := utils.SerializeData(ackObject)
		client.socket.Send(utils.CreateMessageFromByte(ackData));
	}
}

func (client *Client) Emit(eventName string, data interface{}) {
	emitObject := models.GetEmitEventObject(eventName, data, int(client.counter.IncrementAndGet()))
	emitData := utils.SerializeData(emitObject)
	client.socket.Send(utils.CreateMessageFromByte(emitData));
}

func (client *Client) EmitAck(eventName string, data interface{}, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	emitObject := models.GetEmitEventObject(eventName, data, id)
	emitData := utils.SerializeData(emitObject)
	client.putEmitAck(id, eventName, ack)
	client.socket.Send(utils.CreateMessageFromByte(emitData));
}

func (client *Client) Subscribe(channelName string) {
	subscribeObject := models.GetSubscribeEventObject(channelName, int(client.counter.IncrementAndGet()))
	subscribeData := utils.SerializeData(subscribeObject)
	client.socket.Send(utils.CreateMessageFromByte(subscribeData));
}

func (client *Client) SubscribeAck(channelName string, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	subscribeObject := models.GetSubscribeEventObject(channelName, id)
	subscribeData := utils.SerializeData(subscribeObject)
	client.putEmitAck(id, channelName, ack)
	client.socket.Send(utils.CreateMessageFromByte(subscribeData));
}

func (client *Client) Unsubscribe(channelName string) {
	unsubscribeObject := models.GetUnsubscribeEventObject(channelName, int(client.counter.IncrementAndGet()))
	unsubscribeData := utils.SerializeData(unsubscribeObject)
	client.socket.Send(utils.CreateMessageFromByte(unsubscribeData));
}

func (client *Client) UnsubscribeAck(channelName string, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	unsubscribeObject := models.GetUnsubscribeEventObject(channelName, id)
	unsubscribeData := utils.SerializeData(unsubscribeObject)
	client.putEmitAck(id, channelName, ack)
	client.socket.Send(utils.CreateMessageFromByte(unsubscribeData));
}

func (client *Client) Publish(channelName string, data interface{}) {
	publishObject := models.GetPublishEventObject(channelName, data, int(client.counter.IncrementAndGet()))
	publishData := utils.SerializeData(publishObject)
	client.socket.Send(utils.CreateMessageFromByte(publishData));
}

func (client *Client) PublishAck(channelName string, data interface{}, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	publishObject := models.GetPublishEventObject(channelName, data, id)
	publishData := utils.SerializeData(publishObject)
	client.putEmitAck(id, channelName, ack)
	client.socket.Send(utils.CreateMessageFromByte(publishData));
}

func (client *Client) OnChannel(eventName string, ack func(eventName string, data interface{})) {
	client.putOnListener(eventName, ack)
}

func (client *Client) On(eventName string, ack func(eventName string, data interface{})) {
	client.putOnListener(eventName, ack)
}

func (client *Client) OnAck(eventName string, ack func(eventName string, data interface{}, ack func(error interface{}, data interface{}))) {
	client.putOnAckListener(eventName, ack)
}
