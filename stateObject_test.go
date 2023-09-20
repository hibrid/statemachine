package statemachine

import (
	"encoding/json"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
)

// Mock the time for predictable timestamps in tests
var mockTime = time.Date(2023, 9, 7, 0, 0, 0, 0, time.UTC)

func TestNewStateObject(t *testing.T) {
	logger := zaptest.NewLogger(t)
	data := map[string]interface{}{"key": "value"}

	// Create a dummy StateMachine instance
	// The Redis address here is just a placeholder; you can use a mock Redis instance or similar for actual tests
	sm := NewStateMachine("localhost:6379")
	so := NewStateObject(data, sm, logger)

	if so.State != SIMNotActivated {
		t.Errorf("Expected state to be %s but got %s", SIMNotActivated, so.State)
	}

	if val, ok := so.Data["key"]; !ok || val != "value" {
		t.Errorf("Data not set correctly in StateObject")
	}
}

func TestSerialize(t *testing.T) {
	logger := zaptest.NewLogger(t)
	data := map[string]interface{}{"key": "value"}
	// Create a dummy StateMachine instance
	// The Redis address here is just a placeholder; you can use a mock Redis instance or similar for actual tests
	sm := NewStateMachine("localhost:6379")
	so := NewStateObject(data, sm, logger)

	serialized, err := so.Serialize()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(serialized, &result)
	if err != nil {
		t.Fatalf("Deserialization of serialized data failed: %v", err)
	}

	if result["state"] != SIMNotActivated {
		t.Errorf("Serialized state is incorrect")
	}
}

func TestCommitToDisk(t *testing.T) {
	logger := zaptest.NewLogger(t)
	// Mock the CommitFunc
	data := map[string]interface{}{"key": "value"}
	so := NewStateObject(data, nil, logger) // Passing nil as StateMachine as it's not used in this mock
	so.EventID = "test_event"
	so.CommitFunc = func() error {
		return nil // or whatever mock behavior you want
	}

	err := so.CommitToDisk()
	if err != nil {
		t.Fatalf("CommitToDisk failed: %v", err)
	}
}

// TODO: Do we need the TestLogTransition if we are using zap.Logger?

func TestMarkProcessed(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Mock the CommitFunc
	data := map[string]interface{}{"key": "value"}
	so := NewStateObject(data, nil, logger) // Passing nil as StateMachine as it's not used in this mock
	so.EventID = "test_event"
	so.CommitFunc = func() error {
		return nil // or whatever mock behavior you want
	}

	err := so.MarkProcessed()
	if err != nil {
		t.Fatalf("MarkProcessed failed: %v", err)
	}

	if so.State != "Processed" {
		t.Errorf("Expected state to be Processed but got %s", so.State)
	}
}
