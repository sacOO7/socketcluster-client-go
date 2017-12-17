package scclient

import (
	"log"
	_ "golang.org/x/net/websocket"
	"github.com/rgamba/evtwebsocket"
	_ "time"
	"./models"
	"./utils"
	"./parser"
	"./listener"
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
	listener.Listener
}

func New(url string) Client {
	return Client{url: url, counter: utils.AtomicCounter{Counter: 0}, Listener: listener.Init()}
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
					isAuthenticated := GetIsAuthenticated(messageObject)
					if client.onAuthentication != nil {
						client.onAuthentication(*client, isAuthenticated);
					}
				case parser.SETTOKEN:
					token := GetAuthToken(messageObject)
					if client.onSetAuthentication != nil {
						client.onSetAuthentication(*client, token)
					}

				case parser.REMOVETOKEN:
					utils.PrintMessage("got remove token message")
					client.authToken = nil
				case parser.EVENT:
					utils.PrintMessage("got event receive message")
				case parser.ACKRECEIVE:
					client.HandleEmitAck(rid, error, data)
					utils.PrintMessage("got ack receive message")
				case parser.PUBLISH:
					utils.PrintMessage("got publish message")
				}
			}

		},
		// Fires when an error occurs and connection is closed
		OnError: func(err error) {
			if client.onDisconnect != nil {
				client.onDisconnect(*client, err)
			}
		},
		// Ping interval in secs (optional)
		PingIntervalSecs: 5,
		// Ping message to send (optional)
		PingMsg: []byte("PING"),
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

func GetAuthToken(message interface{}) string {
	itemsMap := message.(map[string]interface{})
	data := itemsMap["data"]
	return data.(map[string]interface{})["token"].(string)
}

func GetIsAuthenticated(message interface{}) bool {
	itemsMap := message.(map[string]interface{})
	data := itemsMap["data"]
	return data.(map[string]interface{})["isAuthenticated"].(bool)
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
	client.PutEmitAck(id, eventName, ack)
	client.socket.Send(utils.CreateMessageFromByte(emitData));
}

func (client *Client) Subscribe (channelName string) {
	subscribeObject := models.GetSubscribeEventObject(channelName, int(client.counter.IncrementAndGet()))
	subscribeData := utils.SerializeData(subscribeObject)
	client.socket.Send(utils.CreateMessageFromByte(subscribeData));
}

func (client *Client) SubscribeAck (channelName string, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	subscribeObject := models.GetSubscribeEventObject(channelName, id)
	subscribeData := utils.SerializeData(subscribeObject)
	client.PutEmitAck(id, channelName, ack)
	client.socket.Send(utils.CreateMessageFromByte(subscribeData));
}

func (client *Client) Unsubscribe (channelName string) {
	unsubscribeObject := models.GetUnsubscribeEventObject(channelName, int(client.counter.IncrementAndGet()))
	unsubscribeData := utils.SerializeData(unsubscribeObject)
	client.socket.Send(utils.CreateMessageFromByte(unsubscribeData));
}

func (client *Client) UnsubscribeAck (channelName string, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	unsubscribeObject := models.GetUnsubscribeEventObject(channelName, id)
	unsubscribeData := utils.SerializeData(unsubscribeObject)
	client.PutEmitAck(id, channelName, ack)
	client.socket.Send(utils.CreateMessageFromByte(unsubscribeData));
}

func (client *Client) Publish (channelName string, data interface{}) {
	publishObject := models.GetPublishEventObject(channelName, data, int(client.counter.IncrementAndGet()))
	publishData := utils.SerializeData(publishObject)
	client.socket.Send(utils.CreateMessageFromByte(publishData));
}

func (client *Client) PublishAck (channelName string, data interface{}, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	publishObject := models.GetPublishEventObject(channelName, data, id)
	publishData := utils.SerializeData(publishObject)
	client.PutEmitAck(id, channelName, ack)
	client.socket.Send(utils.CreateMessageFromByte(publishData));
}
