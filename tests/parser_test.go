package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"../parser"
)

func TestShouldReturnPublish(t *testing.T) {
	var expectedParseResult = parser.PUBLISH
	actaulParseResult := parser.Parse(1,1, "#publish")
	assert.Equal(t, expectedParseResult, actaulParseResult, "should be equal")
}

func TestShouldReturnRemoveAuthToken(t *testing.T) {
	var expectedParseResult = parser.REMOVETOKEN
	actaulParseResult := parser.Parse(1,0, "#removeAuthToken")
	assert.Equal(t, expectedParseResult, actaulParseResult, "should be equal")
}

func TestShouldReturnSetAuthToken(t *testing.T) {
	var expectedParseResult = parser.SETTOKEN
	actaulParseResult := parser.Parse(1,0, "#setAuthToken")
	assert.Equal(t, expectedParseResult, actaulParseResult, "should be equal")
}

func TestShouldReturnEvent(t *testing.T) {
	var expectedParseResult = parser.EVENT
	actaulParseResult := parser.Parse(1,0, "chat")
	assert.Equal(t, expectedParseResult, actaulParseResult, "should be equal")
}

func TestShouldReturnIsAuthenticated(t *testing.T) {
	var expectedParseResult = parser.ISAUTHENTICATED
	actaulParseResult := parser.Parse(1,0, nil)
	assert.Equal(t, expectedParseResult, actaulParseResult, "should be equal")
}

func TestShouldReturnAckReceive(t *testing.T) {
	var expectedParseResult = parser.ACKRECEIVE
	actaulParseResult := parser.Parse(12,0, nil)
	assert.Equal(t, expectedParseResult, actaulParseResult, "should be equal")
}


