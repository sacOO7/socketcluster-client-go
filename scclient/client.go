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

var authToken string


func Handle_connection() {
	conn := evtwebsocket.Conn{
		// Fires when the connection is established
		OnConnected: func(w *evtwebsocket.Conn) {
			fmt.Println("Connected!")
			handshake := utils.SerializeData(models.GetHandshakeObject(nil, 1))
			w.Send(utils.CreateMessageFromByte(handshake));
		},
		// Fires when a new message arrives from the server
		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
			fmt.Printf("New message: %s\n", msg)

			if utils.IsEqual("#1", msg) {
				fmt.Println("Got ping message ")
				w.Send(utils.CreateMessageFromString("#2"));
			} else {
				var jsonObject = utils.DeserializeData(msg)
				_, rid, cid, eventname, _ := parser.GetMessageDetails(jsonObject)

				parseresult := parser.Parse(rid, cid, eventname)

				switch parseresult {
					case parser.ISAUTHENTICATED:
						utils.PrintMessage("got is authenticated message")
					case parser.SETTOKEN:
						utils.PrintMessage("got set token message")
					case parser.REMOVETOKEN:
						utils.PrintMessage("got remove token message")
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

	// Connect
	err := conn.Dial("ws://localhost:8000/socketcluster/", "")
	// c.Send([]byte("TEST"), nil)

	if err != nil {
		log.Fatal(err)
	}

	msg := evtwebsocket.Msg{
		Body: []byte("Hello buddy"),
		Callback: func(resp []byte, conn *evtwebsocket.Conn) {
			// This function executes when the server responds
			fmt.Printf("Got response: %s\n", resp)
		},
	}

	conn.Send(msg)

}

