package main

import (
	"./scclient"
	"text/scanner"
	"os"
	"fmt"
	 _"log"
)

func onConnect () {
	fmt.Println("Connected to server")
}

func onDisconnect(err error) {
	fmt.Printf("Error: %s\n", err.Error())
	os.Exit(1)
}

func onConnectError(err error) {
	fmt.Printf("Error: %s\n", err.Error())
	os.Exit(1)
}

func onSetAuthentication (token string) {
	fmt.Println("Auth token received :", token)
}

func  onAuthentication (isAuthenticated bool) {
	fmt.Println("Client authenticated :", isAuthenticated)
}


func main() {
	var reader scanner.Scanner
	client := scclient.New("ws://192.168.0.5:8000/socketcluster/");

	client.SetBasicListener(onConnect, onConnectError, onDisconnect)
	client.SetAuthenticationListener(onSetAuthentication, onAuthentication)
	go client.Connect()

	fmt.Println("Enter any key to terminate the program")
	reader.Init(os.Stdin)
	reader.Next()
	// os.Exit(0)
}
