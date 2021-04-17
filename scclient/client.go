package scclient

import (
	_ "golang.org/x/net/websocket"
	_ "time"
	"github.com/sacOO7/socketcluster-client-go/scclient/models"
	"github.com/sacOO7/socketcluster-client-go/scclient/utils"
	"github.com/sacOO7/socketcluster-client-go/scclient/parser"
	"github.com/sacOO7/gowebsocket"
	"net/http"
	"github.com/sacOO7/go-logger"
)

type Client struct {
	authToken           *string
	url                 string
	counter             utils.AtomicCounter
	Socket              gowebsocket.Socket
	onConnect           func(client Client)
	onConnectError      func(client Client, err error)
	onDisconnect        func(client Client, err error)
	onSetAuthentication func(client Client, token string)
	onAuthentication    func(client Client, isAuthenticated bool)
	ConnectionOptions   gowebsocket.ConnectionOptions
	RequestHeader       http.Header
	Listener
}

func New(url string) Client {

	return Client{
		url:      url,
		counter:  utils.AtomicCounter{Counter: 0},
		Listener: Init()}
}

func (client *Client) IsConnected() bool {
	return client.Socket.IsConnected
}

func (client *Client) EnableLogging() {
	scLogger.SetLevel(logging.TRACE)
}

func (client *Client) GetLogger() logging.Logger {
	return scLogger;
}

func (client *Client) SetAuthToken(token string) {
	client.authToken = &token
}

func (client *Client) GetAuthToken() string {
	return *client.authToken;
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

	client.Socket.OnConnected = func(socket gowebsocket.Socket) {
		client.counter.Reset()
		client.sendHandshake()
		if client.onConnect != nil {
			client.onConnect(*client)
		}
	};
	client.Socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		if err != nil {
			if client.onConnectError != nil {
				client.onConnectError(*client, err)
			}
		}
	};
	client.Socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		scLogger.Info.Printf("%s", message)

		if message == "" {
			client.Socket.SendText("");
		} else {
			var messageObject = utils.DeserializeDataFromString(message)
			data, rid, cid, eventname, error := parser.GetMessageDetails(messageObject)

			parseresult := parser.Parse(rid, cid, eventname)

			switch parseresult {
			case parser.ISAUTHENTICATED:
				isAuthenticated := utils.GetIsAuthenticated(messageObject)
				if client.onAuthentication != nil {
					client.onAuthentication(*client, isAuthenticated);
				}
			case parser.SETTOKEN:
				scLogger.Trace.Println("Set token event received")
				token := utils.GetAuthToken(messageObject)
				if client.onSetAuthentication != nil {
					client.onSetAuthentication(*client, token)
				}

			case parser.REMOVETOKEN:
				scLogger.Trace.Println("Remove token event received")
				client.authToken = nil
			case parser.EVENT:
				scLogger.Trace.Println("Received data for event :: ", eventname)
				if client.hasEventAck(eventname.(string)) {
					client.handleOnAckListener(eventname.(string), data, client.ack(cid))
				} else {
					client.handleOnListener(eventname.(string), data)
				}
			case parser.ACKRECEIVE:
				client.handleEmitAck(rid, error, data)
			case parser.PUBLISH:
				channel := models.GetChannelObject(data)
				scLogger.Trace.Println("Publish event received for channel :: ", channel.Channel)
				client.handleOnListener(channel.Channel, channel.Data)
			}
		}
	};
	client.Socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		if client.onDisconnect != nil {
			client.onDisconnect(*client, err)
		}
		return
	}

}

func (client *Client) Connect() {
	client.Socket = gowebsocket.New(client.url)
	client.registerCallbacks()
	// Connect
	client.Socket.ConnectionOptions = client.ConnectionOptions
	client.Socket.RequestHeader = client.RequestHeader
	client.Socket.Connect()
}

func (client *Client) sendHandshake() {
	handshake := utils.SerializeDataIntoString(models.GetHandshakeObject(client.authToken, int(client.counter.IncrementAndGet())))
	client.Socket.SendText(handshake);
}

func (client *Client) ack(cid int) func(error interface{}, data interface{}) {
	return func(error interface{}, data interface{}) {
		ackObject := models.GetReceiveEventObject(data, error, cid);
		ackData := utils.SerializeDataIntoString(ackObject)
		client.Socket.SendText(ackData);
	}
}

func (client *Client) Emit(eventName string, data interface{}) {
	emitObject := models.GetEmitEventObject(eventName, data, int(client.counter.IncrementAndGet()))
	emitData := utils.SerializeDataIntoString(emitObject)
	client.Socket.SendText(emitData);
}

func (client *Client) EmitAck(eventName string, data interface{}, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	emitObject := models.GetEmitEventObject(eventName, data, id)
	emitData := utils.SerializeDataIntoString(emitObject)
	client.putEmitAck(id, eventName, ack)
	client.Socket.SendText(emitData);
}

func (client *Client) Subscribe(channelName string) {
	subscribeObject := models.GetSubscribeEventObject(channelName, int(client.counter.IncrementAndGet()))
	subscribeData := utils.SerializeDataIntoString(subscribeObject)
	client.Socket.SendText(subscribeData);
}

func (client *Client) SubscribeAck(channelName string, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	subscribeObject := models.GetSubscribeEventObject(channelName, id)
	subscribeData := utils.SerializeDataIntoString(subscribeObject)
	client.putEmitAck(id, channelName, ack)
	client.Socket.SendText(subscribeData);
}

func (client *Client) Unsubscribe(channelName string) {
	unsubscribeObject := models.GetUnsubscribeEventObject(channelName, int(client.counter.IncrementAndGet()))
	unsubscribeData := utils.SerializeDataIntoString(unsubscribeObject)
	client.Socket.SendText(unsubscribeData);
}

func (client *Client) UnsubscribeAck(channelName string, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	unsubscribeObject := models.GetUnsubscribeEventObject(channelName, id)
	unsubscribeData := utils.SerializeDataIntoString(unsubscribeObject)
	client.putEmitAck(id, channelName, ack)
	client.Socket.SendText(unsubscribeData);
}

func (client *Client) Publish(channelName string, data interface{}) {
	publishObject := models.GetPublishEventObject(channelName, data, int(client.counter.IncrementAndGet()))
	publishData := utils.SerializeDataIntoString(publishObject)
	client.Socket.SendText(publishData);
}

func (client *Client) PublishAck(channelName string, data interface{}, ack func(eventName string, error interface{}, data interface{})) {
	id := int(client.counter.IncrementAndGet())
	publishObject := models.GetPublishEventObject(channelName, data, id)
	publishData := utils.SerializeDataIntoString(publishObject)
	client.putEmitAck(id, channelName, ack)
	client.Socket.SendText(publishData);
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

func (client *Client) Disconnect() {
	client.Socket.Close()
}
