package structs

import "testing"

// TestNewResponse 验证单题判题响应的构造。
func TestNewResponse(t *testing.T) {
	response := NewResponse(true, "explanation")
	if !response.Correct || response.Explanation != "explanation" {
		t.Fatalf("unexpected response: %+v", response)
	}
}
