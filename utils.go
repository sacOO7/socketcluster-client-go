package main
import (
	"fmt"
	"github.com/rgamba/evtwebsocket"
    "encoding/json"
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

func CreateMessageFromString(message string) evtwebsocket.Msg{
	return evtwebsocket.Msg{
    Body: []byte(message),
	}
}

func CreateMessageFromByte(message [] byte) evtwebsocket.Msg{
    return evtwebsocket.Msg{
    Body: message,
    }
}

func SerializeData(data interface {}) [] byte {
    b, _ := json.Marshal(data)
    return b;
}
