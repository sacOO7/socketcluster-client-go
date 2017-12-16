package scclient

import (
	"fmt"
	_ "golang.org/x/net/websocket"
	"github.com/rgamba/evtwebsocket"
	"log"
	_ "time"
	"os"
	"./models"
	"./utils"
	"./parser"
)

type Client struct {
	authToken *string
	url       string
	counter   utils.AtomicCounter
	socket    *evtwebsocket.Conn
}

func NewClient(url string) Client {
	return Client{url: url, counter: utils.AtomicCounter{Counter: 0}}
}

func (client *Client) registerCallbacks() {
	client.socket = & evtwebsocket.Conn{
		// Fires when the connection is established
		OnConnected: func(w *evtwebsocket.Conn) {
			fmt.Println("Connected!")
			client.sendHandshake()
		},
		// Fires when a new message arrives from the server
		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
			fmt.Printf("New message: %s\n", msg)

			if utils.IsEqual("#1", msg) {
				fmt.Println("Got ping message ")
				w.Send(utils.CreateMessageFromString("#2"));
			} else {
				var messageObject = utils.DeserializeData(msg)
				_, rid, cid, eventname, _ := parser.GetMessageDetails(messageObject)

				parseresult := parser.Parse(rid, cid, eventname)

				switch parseresult {
				case parser.ISAUTHENTICATED:
					isAuthenticated := GetIsAuthenticated(messageObject)
					fmt.Println("Is authenticated is ", isAuthenticated)
					utils.PrintMessage("got is authenticated message")
				case parser.SETTOKEN:
					token := GetAuthToken(messageObject)
					fmt.Println("token is ", token)
					utils.PrintMessage("got set token message")
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
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)

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
		log.Fatal(err)
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
