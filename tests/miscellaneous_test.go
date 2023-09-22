package tests

import (
	"testing"

	"github.com/sacOO7/socketcluster-client-go/scclient/models"
	"github.com/sacOO7/socketcluster-client-go/scclient/utils"

	"github.com/stretchr/testify/assert"
)

func TestShouldCheckEqual(t *testing.T) {
	expected := []byte("mystring")
	assert.True(t, utils.IsEqual("mystring", expected), "String and byte [] should be equal")

}

func TestShouldSerializeData(t *testing.T) {

	emitEvent := models.EmitEvent{Cid: 2, Data: "My sample data", Event: "chat"}
	expectedData := "{\"event\":\"chat\",\"data\":\"My sample data\",\"cid\":2}"

	assert.Equal(t, expectedData, string(utils.SerializeData(emitEvent)))
}
