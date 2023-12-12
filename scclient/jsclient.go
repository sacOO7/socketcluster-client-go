package scclient

import (
	"github.com/sacOO7/socketcluster-client-go/scclient/models"
	"github.com/sacOO7/socketcluster-client-go/scclient/utils"
)

func (client *Client) Transmit(eventName string, data interface{}) {
	transmitObject := models.GetTransmitEventObject(eventName, data)
	transmitData := utils.SerializeDataIntoString(transmitObject)
	client.socket.SendText(transmitData)
}