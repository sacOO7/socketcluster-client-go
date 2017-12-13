package main
import (
	"fmt"
	_ "golang.org/x/net/websocket"
	"github.com/rgamba/evtwebsocket"
	"log"
	_ "time"
	"os"
	"./models"
	"./utils"
)

func handle_connection() {
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
  err := conn.Dial("ws://localhost:8000/socketcluster/","")
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


func main() {
	var i int
	go handle_connection()

	utils.PrintMe("Enter any key to terminate the program")
	fmt.Scan(&i)
	// os.Exit(0)
}

