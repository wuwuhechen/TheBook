package structs

import "testing"

func TestNewResponse(t *testing.T) {
	response := NewResponse(true, "explanation")
	if !response.Correct || response.Explanation != "explanation" {
		t.Fatalf("unexpected response: %+v", response)
	}
}
