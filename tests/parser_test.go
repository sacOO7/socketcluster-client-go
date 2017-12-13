package tests

import (
	"testing"
)

func TestReverseToReturnReversedInputString(t *testing.T) {
	actualResult := "olleH"
	var expectedResult = "olleH"

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}


