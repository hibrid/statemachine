package statemachine

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap/zaptest"
)

// Dummy Redis client for testing
var mockRedisClient *redis.Client

// Mock redis functions
func mockRedisExists(eventID string) bool {
	return false // Or whatever mock behavior you want
}

func mockRedisSet(eventID string, value bool) {}

// Initialize dummy Redis client for testing
func init() {
	mockRedisClient = redis.NewClient(&redis.Options{})
}

func TestSetHandlerConfig(t *testing.T) {
	sm := &StateMachine{}
	config := HandlerConfig{
		MarkProcessed: true,
		CheckEventID:  true,
	}
	sm.SetHandlerConfig(config)
	if sm.config != config {
		t.Fatalf("HandlerConfig not set correctly")
	}
}

func TestInitializeRedis(t *testing.T) {
	rdb := InitializeRedis("localhost:6379")
	if rdb == nil {
		t.Fatalf("Redis client not initialized")
	}
}

func TestNewStateMachine(t *testing.T) {
	sm := NewStateMachine("localhost:6379")
	if sm.redisClient == nil {
		t.Fatalf("Redis client not initialized in StateMachine")
	}
}

func TestStateMachineLog(t *testing.T) {
	sm := &StateMachine{
		LogTransitions: true,
	}
	sm.Log("Test message") // This should not produce an error
}

func TestStateMachineLogErr(t *testing.T) {
	sm := &StateMachine{
		LogTransitions: true,
		DebugLogging:   false,
	}
	sm.LogErr(errors.New("Test error")) // This should not produce an error

	sm.DebugLogging = true
	sm.LogErr(errors.New("Test error with debug")) // This should not produce an error
}

func TestGetDefaultHandlers(t *testing.T) {
	sm := &StateMachine{
		config: defaultConfig,
	}
	handlers := sm.GetDefaultHandlers()
	if len(handlers) == 0 {
		t.Fatalf("No default handlers returned")
	}
}

func TestSetConfig(t *testing.T) {
	sm := &StateMachine{}
	config := Config{
		DB: nil, // Mock other properties as needed
	}
	sm.SetConfig(config)
	if sm.db != nil {
		t.Fatalf("DB should not be set in StateMachine")
	}
}

func TestSetDB(t *testing.T) {
	sm := &StateMachine{}
	db := &sql.DB{}
	sm.SetDB(db)
	if sm.db != db {
		t.Fatalf("DB not set correctly in StateMachine")
	}
}

func TestStateMachineEmitEvent(t *testing.T) {
	sm := &StateMachine{}
	sm.EmitEvent("Test event") // This should not produce an error
}

func TestGetHandlerType(t *testing.T) {
	handler := &CheckEventIDHandler{}
	handlerType := getHandlerType(handler)
	if handlerType != "CheckEventIDHandler" {
		t.Fatalf("Expected handler type to be CheckEventIDHandler but got %s", handlerType)
	}
}

func TestRegisterTransition(t *testing.T) {
	sm := NewStateMachine("localhost:6379")
	customHandler := &CustomCheckEventIDHandler{}
	sm.RegisterTransition("FromState", "ToState", customHandler)
	_, exists := sm.transitions["FromState->ToState"]
	if !exists {
		t.Fatalf("Transition not registered correctly")
	}
}

func TestGenerateSignature(t *testing.T) {
	signature := generateSignature("Test event")
	if signature == "" {
		t.Fatalf("Signature generation failed")
	}
}

func TestTransitionTo(t *testing.T) {
	logger := zaptest.NewLogger(t)

	sm := NewStateMachine("localhost:6379")
	data := map[string]interface{}{"key": "value"}
	so := NewStateObject(data, sm, logger)
	err := so.TransitionTo(sm, "ToState")
	if err == nil {
		t.Fatalf("Transition should fail due to no registered transitions")
	}
}

func TestStateMachineIsProcessed(t *testing.T) {
	sm := NewStateMachine("localhost:6379")
	sm.redisClient = mockRedisClient
	if sm.isProcessed("test_event") {
		t.Fatalf("Event should not be processed")
	}
}

func TestStateMachineMarkProcessed(t *testing.T) {
	sm := NewStateMachine("localhost:6379")
	sm.redisClient = mockRedisClient
	sm.markProcessed("test_event") // This should not produce an error
}
