package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"../scclient/parser"
	"encoding/json"
)

func TestShouldReturnPublish(t *testing.T) {
	var expectedParseResult = parser.PUBLISH
	actaulParseResult := parser.Parse(1, 1, "#publish")
	assert.Equal(t, expectedParseResult, actaulParseResult, "should return publish")
}

func TestShouldReturnRemoveAuthToken(t *testing.T) {
	var expectedParseResult = parser.REMOVETOKEN
	actaulParseResult := parser.Parse(1, 0, "#removeAuthToken")
	assert.Equal(t, expectedParseResult, actaulParseResult, "should return remove auth token")
}

func TestShouldReturnSetAuthToken(t *testing.T) {
	var expectedParseResult = parser.SETTOKEN
	actaulParseResult := parser.Parse(1, 0, "#setAuthToken")
	assert.Equal(t, expectedParseResult, actaulParseResult, "should return set auth token")
}

func TestShouldReturnEvent(t *testing.T) {
	var expectedParseResult = parser.EVENT
	actaulParseResult := parser.Parse(1, 0, "chat")
	assert.Equal(t, expectedParseResult, actaulParseResult, "should return custom event")
}

func TestShouldReturnIsAuthenticated(t *testing.T) {
	var expectedParseResult = parser.ISAUTHENTICATED
	actaulParseResult := parser.Parse(1, 0, nil)
	assert.Equal(t, expectedParseResult, actaulParseResult, "should return is authenticated event")
}

func TestShouldReturnAckReceive(t *testing.T) {
	var expectedParseResult = parser.ACKRECEIVE
	actaulParseResult := parser.Parse(12, 0, nil)
	assert.Equal(t, expectedParseResult, actaulParseResult, "should return ack receive event")
}

func TestShouldReturnMessageDetails(t *testing.T) {
	message := "{\"event\":\"#removeAuthToken\",\"data\":\"This is a sample data\",\"cid\":1, \"rid\":2, \"error\":\"This is a sample error\"}"
	var jsonObject interface{}
	json.Unmarshal([] byte (message), &jsonObject)
	data, rid, cid, eventname, error := parser.GetMessageDetails(jsonObject)
	assert.Equal(t, "This is a sample data", data)
	assert.Equal(t, 2, rid)
	assert.Equal(t, 1, cid)
	assert.Equal(t, "#removeAuthToken", eventname)
	assert.Equal(t, "This is a sample error", error)

}
