package statemachine

import (
	"database/sql"
	"time"

	"github.com/go-redis/redis/v8"
)

type StateMachine struct {
	transitions    map[string]StateTransition
	config         HandlerConfig
	LogTransitions bool
	DebugLogging   bool
	db             *sql.DB
	eventEmitter   interface{}
	kafkaConn      interface{}
	natsConn       interface{}
	rabbitMQConn   interface{}
	redisClient    *redis.Client
}

type Config struct {
	DB           *sql.DB
	KafkaConn    interface{}
	NatsConn     interface{}
	RabbitMQConn interface{}
}

type StateTransition struct {
	From  string
	To    string
	Chain Handler
}

type Handler interface {
	Handle(*StateObject, string) bool
	SetNext(Handler)
	Rollback(*StateObject, string) bool
	Next() Handler
}

type HandlerConfig struct {
	CheckEventID   bool
	CheckProcessed bool
	Telemetry      bool
	Alerting       bool
	MarkProcessed  bool
}

type StateTransitionLog struct {
	EventID   string
	FromState string
	ToState   string
	Timestamp time.Time
}
