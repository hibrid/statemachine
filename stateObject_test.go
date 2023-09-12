package statemachine

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// Mock the time for predictable timestamps in tests
var mockTime = time.Date(2023, 9, 7, 0, 0, 0, 0, time.UTC)

func TestNewStateObject(t *testing.T) {
	data := map[string]interface{}{"key": "value"}

	// Create a dummy StateMachine instance
	// The Redis address here is just a placeholder; you can use a mock Redis instance or similar for actual tests
	sm := NewStateMachine("localhost:6379")
	so := NewStateObject(data, sm)

	if so.State != SIMNotActivated {
		t.Errorf("Expected state to be %s but got %s", SIMNotActivated, so.State)
	}

	if val, ok := so.Data["key"]; !ok || val != "value" {
		t.Errorf("Data not set correctly in StateObject")
	}
}

func TestSerialize(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	// Create a dummy StateMachine instance
	// The Redis address here is just a placeholder; you can use a mock Redis instance or similar for actual tests
	sm := NewStateMachine("localhost:6379")
	so := NewStateObject(data, sm)

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
	// Mock the CommitFunc
	data := map[string]interface{}{"key": "value"}
	so := NewStateObject(data, nil) // Passing nil as StateMachine as it's not used in this mock
	so.EventID = "test_event"
	so.CommitFunc = func() error {
		return nil // or whatever mock behavior you want
	}

	err := so.CommitToDisk()
	if err != nil {
		t.Fatalf("CommitToDisk failed: %v", err)
	}
}

func TestLogTransition(t *testing.T) {
	// Mock fmt.Printf to capture printed logs
	var capturedLog string
	mockLogger := func(format string, a ...interface{}) (n int, err error) {
		capturedLog = fmt.Sprintf(format, a...)
		return 0, nil
	}

	// Mock time.Now
	oldNow := nowFunc
	defer func() { nowFunc = oldNow }()
	nowFunc = func() time.Time { return mockTime }

	data := map[string]interface{}{"key": "value"}
	sm := NewStateMachine("localhost:6379")
	so := NewStateObject(data, sm)
	so.Logger = mockLogger
	so.EventID = "test_event"

	sm.DebugLogging = true

	so.LogTransition("fromState", "toState", sm)

	expectedLog := fmt.Sprintf("{timestamp: %s, event_id: %s, from_state: fromState, to_state: toState, func: git.mena.technology/statemachine.TestLogTransition, file: stateObject_test.go, line: 89}\n", mockTime, so.EventID)
	if capturedLog != expectedLog {
		t.Errorf("Expected log to be %s but got %s", expectedLog, capturedLog)
	}
}

func TestMarkProcessed(t *testing.T) {
	// Mock the CommitFunc
	data := map[string]interface{}{"key": "value"}
	so := NewStateObject(data, nil) // Passing nil as StateMachine as it's not used in this mock
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
