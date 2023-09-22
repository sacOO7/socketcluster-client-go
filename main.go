package main

import (
	"fmt"
	"os"
	"text/scanner"

	"github.com/sacOO7/socketcluster-client-go/scclient"
)

func onConnect(client scclient.Client) {
	fmt.Println("Connected to server")
}

func onDisconnect(client scclient.Client, err error) {
	fmt.Printf("Error: %s\n", err.Error())
	os.Exit(1)
}

func onConnectError(client scclient.Client, err error) {
	fmt.Printf("Error: %s\n", err.Error())
	os.Exit(1)
}

func onSetAuthentication(client scclient.Client, token string) {
	fmt.Println("Auth token received :", token)

}

func onAuthentication(client scclient.Client, isAuthenticated bool) {
	fmt.Println("Client authenticated :", isAuthenticated)
}

func main() {
	var reader scanner.Scanner
	client := scclient.New("ws://localhost:8000/socketcluster/")
	client.SetBasicListener(onConnect, onConnectError, onDisconnect)
	client.SetAuthenticationListener(onSetAuthentication, onAuthentication)
	client.EnableLogging()
	go client.Connect()

	fmt.Println("Enter any key to terminate the program")
	reader.Init(os.Stdin)
	reader.Next()
	// os.Exit(0)
}

func start(client scclient.Client) {
	// start writing your code from here
}
