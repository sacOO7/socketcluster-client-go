package main

import (
	"fmt"
)

func parse (rid int, cid int, event interface {}) {
	if event != nil {
		if event == "#publish" {
			fmt.Println("This is my first publish")

		} else if event == "#removeAuthToken" {
			fmt.Println("This is a removeAuthToken")

		} else if event == "#setAuthToken" {
			fmt.Println("This is a setAuthToken")

		} else {
			fmt.Println("THis is event got called from server")
		}
	} else if rid == 1 {
			fmt.Println("is authenticated got called")
	} else {
			fmt.Println("got ack")
	}
}
