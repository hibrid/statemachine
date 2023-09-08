package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	SIMNotActivated     = "SIMNotActivated"
	SIMActivated        = "SIMActivated"
	PhoneNumberRecycled = "PhoneNumberRecycled"
	SIMDeactivated      = "SIMDeactivated"
	BillingPaid         = "BillingPaid"
	BillingFailed       = "BillingFailed"
	ManualReview        = "ManualReview"
)

var defaultConfig = HandlerConfig{
	CheckEventID:   true,
	CheckProcessed: true,
	Telemetry:      true,
	Alerting:       true,
	MarkProcessed:  true,
}

func (sm *StateMachine) SetHandlerConfig(config HandlerConfig) {
	sm.config = config
}

// Initialize Redis client
func InitializeRedis(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr, // use default Addr
		Password: "",   // no password
		DB:       0,    // use default DB
	})
	return rdb
}

func NewStateMachine(redisAddr string) *StateMachine {
	return &StateMachine{
		transitions: make(map[string]StateTransition),
		config:      defaultConfig,
		redisClient: InitializeRedis(redisAddr),
	}
}

// Adding a Log function to the StateMachine
func (sm *StateMachine) Log(messages ...interface{}) {
	if sm.LogTransitions {
		fmt.Println(messages...)
	}
}

// Adding a LogErr function to the StateMachine
func (sm *StateMachine) LogErr(err error) {
	if sm.LogTransitions {
		sm.Log("Error:", err)
		if sm.DebugLogging {
			_, file, line, _ := runtime.Caller(1)
			fmt.Println("Debug:", file, line)
		}
	}
}

func (sm *StateMachine) GetDefaultHandlers() []Handler {
	var defaultHandlers []Handler
	if sm.config.CheckEventID {
		defaultHandlers = append(defaultHandlers, &CheckEventIDHandler{})
	}
	if sm.config.CheckProcessed {
		defaultHandlers = append(defaultHandlers, &CheckProcessedHandler{
			stateMachine: sm,
		})
	}
	if sm.config.Telemetry {
		defaultHandlers = append(defaultHandlers, &TelemetryHandler{telemetry: sm.config.Telemetry})
	}
	if sm.config.Alerting {
		defaultHandlers = append(defaultHandlers, &AlertingHandler{alerting: sm.config.Alerting})
	}
	return defaultHandlers
}

func (sm *StateMachine) SetConfig(config Config) {
	if config.DB != nil {
		sm.db = config.DB
	}

	if config.KafkaConn != nil {
		if sm.eventEmitter != nil {
			log.Println("Warning: Overwriting existing event emitter with Kafka connection.")
		}
		sm.eventEmitter = config.KafkaConn
		sm.kafkaConn = config.KafkaConn
	}

	if config.NatsConn != nil {
		if sm.eventEmitter != nil {
			log.Println("Warning: Overwriting existing event emitter with NATS connection.")
		}
		sm.eventEmitter = config.NatsConn
		sm.natsConn = config.NatsConn
	}

	if config.RabbitMQConn != nil {
		if sm.eventEmitter != nil {
			log.Println("Warning: Overwriting existing event emitter with RabbitMQ connection.")
		}
		sm.eventEmitter = config.RabbitMQConn
		sm.rabbitMQConn = config.RabbitMQConn
	}
}

func (sm *StateMachine) SetDB(db *sql.DB) {
	sm.db = db
}

func (sm *StateMachine) EmitEvent(event interface{}) {
	if sm.eventEmitter == nil {
		log.Println("Warning: No event emitter set. Unable to emit event.")
		return
	}
	// Emit event using the appropriate connection based on the event type or configuration
	// This is a placeholder. Implement the actual event emitting logic here.
}

func getHandlerType(h Handler) string {
	switch h.(type) {
	case *CheckEventIDHandler, *CustomCheckEventIDHandler:
		return "CheckEventIDHandler"
	case *CheckProcessedHandler, *CustomCheckProcessedHandler:
		return "CheckProcessedHandler"
	case *TelemetryHandler, *CustomTelemetryHandler:
		return "TelemetryHandler"
	case *AlertingHandler, *CustomAlertingHandler:
		return "AlertingHandler"
	case *MarkProcessedHandler, *CustomMarkProcessedHandler:
		return "MarkProcessedHandler"
	// ... other cases can be added as needed ...
	default:
		return ""
	}
}

func (sm *StateMachine) RegisterTransition(from, to string, customHandlers ...Handler) {
	allHandlers := sm.GetDefaultHandlers()

	// Map to track custom handlers
	customHandlerMap := make(map[string]Handler)
	handlersToRemove := []int{}
	for i, handler := range customHandlers {
		handlerType := getHandlerType(handler)
		if handlerType != "" {
			customHandlerMap[handlerType] = handler
			handlersToRemove = append(handlersToRemove, i)
		}
	}

	// Remove the custom handlers that were added to the map from the customHandlers slice
	for i := len(handlersToRemove) - 1; i >= 0; i-- {
		index := handlersToRemove[i]
		customHandlers = append(customHandlers[:index], customHandlers[index+1:]...)
	}

	// Replace default handlers with custom ones where provided
	for i, defaultHandler := range allHandlers {
		handlerType := getHandlerType(defaultHandler)
		if customHandler, ok := customHandlerMap[handlerType]; ok {
			allHandlers[i] = customHandler
		}
	}

	// Confirm that all handlers are chained
	for i := 0; i < len(allHandlers)-1; i++ {
		allHandlers[i].SetNext(allHandlers[i+1])
	}

	lastHandler := allHandlers[len(allHandlers)-1]
	for _, handler := range customHandlers {
		lastHandler.SetNext(handler)
		lastHandler = handler
	}

	if sm.config.MarkProcessed {
		markProcessedHandler := &MarkProcessedHandler{
			stateMachine: sm,
		}
		lastHandler.SetNext(markProcessedHandler)
	}

	sm.transitions[from+"->"+to] = StateTransition{
		From:  from,
		To:    to,
		Chain: allHandlers[0],
	}
}

func generateSignature(eventContent string) string {
	timestamp := time.Now().Format("2006-01-02")
	hash := sha256.Sum256([]byte(eventContent + timestamp))
	return hex.EncodeToString(hash[:])
}

func (so *StateObject) TransitionTo(sm *StateMachine, state string) error {
	sm.Log("Starting transition from", so.State, "to", state)

	transition, exists := sm.transitions[so.State+"->"+state]
	if !exists {
		return errors.New("invalid transition from " + so.State + " to " + state)
	}

	handler := transition.Chain
	var executedHandlers []Handler
	for handler != nil {
		success := handler.Handle(so, state)
		if !success {
			// Rollback
			// Log failure in the handler chain
			sm.LogErr(fmt.Errorf("Handler %s failed for eventID %s", getHandlerType(handler), so.EventID))
			for i := len(executedHandlers) - 1; i >= 0; i-- {
				if !executedHandlers[i].Rollback(so, state) {
					// Log failure in the handler chain
					sm.LogErr(fmt.Errorf("Handler %s failed to rollback for eventID %s", getHandlerType(handler), so.EventID))
					so.State = ManualReview
					return errors.New("failed to rollback, moving to manual review")
				}
			}
			return errors.New("transition failed")
		}
		executedHandlers = append(executedHandlers, handler)
		handler = handler.Next()
	}

	// Log the conclusion of the transition
	sm.Log("Successfully concluded transition from", so.State, "to", state)
	so.State = state
	return nil
}

func (sm *StateMachine) isProcessed(eventID string) bool {
	exists, err := sm.redisClient.Exists(context.Background(), eventID).Result()
	if err != nil {
		// Handle error
		return false
	}
	return exists == 1
}

func (sm *StateMachine) markProcessed(eventID string) {
	err := sm.redisClient.Set(context.Background(), eventID, true, 0).Err()
	if err != nil {
		// Handle error
	}
}
