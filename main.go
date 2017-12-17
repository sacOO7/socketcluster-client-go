package main

import (
	"./scclient"
	"text/scanner"
	"os"
	"fmt"
	_ "log"
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
	go start(client)
}

func main() {
	var reader scanner.Scanner
	client := scclient.New("ws://192.168.100.11:8000/socketcluster/");
	client.SetBasicListener(onConnect, onConnectError, onDisconnect)
	client.SetAuthenticationListener(onSetAuthentication, onAuthentication)
	go client.Connect()

	fmt.Println("Enter any key to terminate the program")
	reader.Init(os.Stdin)
	reader.Next()
	// os.Exit(0)
}

func start(client scclient.Client) {

}
