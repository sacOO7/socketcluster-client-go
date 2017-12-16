package scclient

import (
	"log"
	_ "golang.org/x/net/websocket"
	"github.com/rgamba/evtwebsocket"
	_ "time"
	"./models"
	"./utils"
	"./parser"
)

type Client struct {
	authToken *string
	url       string
	counter   utils.AtomicCounter
	socket    *evtwebsocket.Conn
	onConnect func()
	onConnectError func(err error)
	onDisconnect func(err error)
	onSetAuthentication func(token string)
	onAuthentication func(isAuthenticated bool)
}

func New(url string) Client {
	return Client{url: url, counter: utils.AtomicCounter{Counter: 0}}
}

func (client *Client) SetBasicListener(onConnect func(), onConnectError func(err error), onDisconnect func(err error)) {
	client.onConnect = onConnect
	client.onConnectError = onConnectError
	client.onDisconnect = onDisconnect
}

func (client *Client) SetAuthenticationListener(onSetAuthentication func(token string), onAuthentication func(isAuthenticated bool)) {
	client.onSetAuthentication = onSetAuthentication
	client.onAuthentication = onAuthentication
}

func (client *Client) registerCallbacks() {
	client.socket = & evtwebsocket.Conn{
		// Fires when the connection is established
		OnConnected: func(w *evtwebsocket.Conn) {
			if client.onConnect != nil {
				client.onConnect()
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
				_, rid, cid, eventname, _ := parser.GetMessageDetails(messageObject)

				parseresult := parser.Parse(rid, cid, eventname)

				switch parseresult {
				case parser.ISAUTHENTICATED:
					isAuthenticated := GetIsAuthenticated(messageObject)
					if client.onAuthentication != nil {
						client.onAuthentication(isAuthenticated);
					}
				case parser.SETTOKEN:
					token := GetAuthToken(messageObject)
					if client.onSetAuthentication != nil {
						client.onSetAuthentication(token)
					}

				case parser.REMOVETOKEN:
					utils.PrintMessage("got remove token message")
					client.authToken = nil
				case parser.EVENT:
					utils.PrintMessage("got event receive message")
				case parser.ACKRECEIVE:
					utils.PrintMessage("got ack receive message")
				case parser.PUBLISH:
					utils.PrintMessage("got publish message")
				}
			}

		},
		// Fires when an error occurs and connection is closed
		OnError: func(err error) {
			if client.onDisconnect != nil {
				client.onDisconnect(err)
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
			client.onConnectError(err)
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
