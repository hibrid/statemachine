package statemachine

import (
	"testing"
)

func TestDeserialize(t *testing.T) {
	// Setup test data
	data := []byte(`{"State":"Active","Data":{"key":"value"}}`)
	expectedState := "Active"
	expectedData := map[string]interface{}{
		"key": "value",
	}

	// Call Deserialize function
	so, err := Deserialize(data)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	// Validate output
	if so.State != expectedState {
		t.Errorf("Expected state %s but got %s", expectedState, so.State)
	}

	if d, ok := so.Data["key"]; !ok || d != expectedData["key"] {
		t.Errorf("Expected data %v but got %v", expectedData, so.Data)
	}
}

func TestDeserialize_InvalidInput(t *testing.T) {
	// Setup test data
	data := []byte(`{"State":"Active", "Data":}`) // Invalid JSON

	// Call Deserialize function
	_, err := Deserialize(data)
	if err == nil {
		t.Fatalf("Expected error for invalid input but got nil")
	}
}

type MockHandler struct {
	next Handler
}

func (h *MockHandler) Handle(so *StateObject, state string) bool {
	return true
}

func (h *MockHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *MockHandler) Rollback(so *StateObject, state string) bool {
	return true
}

func (h *MockHandler) Next() Handler {
	return h.next
}

func TestCreateHandlerChain(t *testing.T) {
	handler1 := &MockHandler{}
	handler2 := &MockHandler{}
	handler3 := &MockHandler{}

	chainStart := CreateHandlerChain(handler1, handler2, handler3)

	// Check if the start of the chain is handler1
	if start, ok := chainStart.(*MockHandler); !ok || start != handler1 {
		t.Fatalf("Expected chain start to be handler1 but got %v", chainStart)
	}

	// Check if handler1 points to handler2
	if handler1.next != handler2 {
		t.Errorf("Expected handler1's next to be handler2 but got %v", handler1.next)
	}

	// Check if handler2 points to handler3
	if handler2.next != handler3 {
		t.Errorf("Expected handler2's next to be handler3 but got %v", handler2.next)
	}

	// Check if handler3 points to nil (end of the chain)
	if handler3.next != nil {
		t.Errorf("Expected handler3's next to be nil but got %v", handler3.next)
	}
}
