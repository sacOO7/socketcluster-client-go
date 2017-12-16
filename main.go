package main

import (
	"./scclient"
	"./scclient/utils"
	"text/scanner"
	"os"
)

func main() {
	var reader scanner.Scanner
	client := scclient.NewClient("ws://localhost:8000/socketcluster/");
	go client.Connect()

	utils.PrintMessage("Enter any key to terminate the program")
	reader.Init(os.Stdin)
	reader.Next()
	// os.Exit(0)
}
