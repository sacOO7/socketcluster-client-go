package main
import (
	"fmt"
	"github.com/rgamba/evtwebsocket"
)

func PrintMe(message string) {
	fmt.Println(message)
}

func Equal(s string, b []byte) bool {
    if len(s) != len(b) {
        return false
    }
    for i, x := range b {
        if x != s[i] {
            return false
        }
    }
    return true
}

func CreateMessage(message string) evtwebsocket.Msg{
	return evtwebsocket.Msg{
    Body: []byte(message),
	}
}

