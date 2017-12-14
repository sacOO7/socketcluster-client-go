package main

import (
	"fmt"
	"./scclient"
	"./scclient/utils"
)

func main() {
	var i int
	go scclient.Handle_connection()

	utils.PrintMessage("Enter any key to terminate the program")
	fmt.Scan(&i)
	// os.Exit(0)
}
