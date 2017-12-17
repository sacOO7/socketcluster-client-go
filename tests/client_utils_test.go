package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"github.com/sacOO7/socketcluster-client-go/scclient/utils"
)

func TestShouldReturnAuthToken(t *testing.T) {
	event := "{\"event\":\"#setAuthToken\",\"data\": {\"token\":\"234234\"},\"cid\": 2}"
	var jsonObject interface{}
	json.Unmarshal([] byte (event), &jsonObject)
	actualtoken := utils.GetAuthToken(jsonObject)
	assert.Equal(t, "234234", actualtoken)
}

func TestShouldReturnAuthenticationFlag(t *testing.T) {
	event := "{\"rid\":1,\"data\":{\"id\":\"nhI9Ry88h_XpLHwEAAAF\",\"isAuthenticated\":false,\"pingTimeout\":20000}}"
	var jsonObject interface{}
	json.Unmarshal([] byte (event), &jsonObject)
	actualAuthenticationFlag := utils.GetIsAuthenticated(jsonObject)
	assert.Equal(t, false, actualAuthenticationFlag)

}
