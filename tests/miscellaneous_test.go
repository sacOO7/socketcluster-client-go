package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"../scclient/utils"
	"../scclient/models"
)

func TestShouldCheckEqual(t *testing.T) {
	expected := [] byte("mystring")
	assert.True(t, utils.IsEqual("mystring", expected), "String and byte [] should be equal")

}

func TestShouldSerializeData(t *testing.T) {

	emitEvent := models.EmitEvent{Cid: 2, Data: "My sample data", Event: "chat"}
	expectedData := "{\"event\":\"chat\",\"data\":\"My sample data\",\"cid\":2}"

	assert.Equal(t, expectedData, string(utils.SerializeData(emitEvent)))
}
