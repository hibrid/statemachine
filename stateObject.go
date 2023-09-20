package statemachine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
)

type StateObject struct {
	Data       map[string]interface{} `json:"data"`
	State      string                 `json:"state"`
	EventID    string                 `json:"eventID"`
	Logger     *zap.Logger
	CommitFunc func() error `json:"-"`
}

func NewStateObjectFromStruct(data interface{}, sm *StateMachine, logger *zap.Logger) *StateObject {
	var state = &StateObject{
		State:  SIMNotActivated, // or some other default state
		Logger: logger,
	}
	err := state.EncodeObjectToData(data)
	if err != nil {
		panic(err)
	}
	state.CommitFunc = func() error {
		return state.actualCommitToDisk(sm)
	}
	return state
}

func NewStateObject(data map[string]interface{}, sm *StateMachine, logger *zap.Logger) *StateObject {
	var state = &StateObject{
		Data:   data,
		State:  SIMNotActivated, // or some other default state
		Logger: logger,
	}
	state.CommitFunc = func() error {
		return state.actualCommitToDisk(sm)
	}
	return state
}

func (so *StateObject) CommitToDisk() error {
	return so.CommitFunc()
}

func (so *StateObject) Serialize() ([]byte, error) {
	return json.Marshal(so)
}

var writeFileFunc = os.WriteFile

func (so *StateObject) actualCommitToDisk(sm *StateMachine) error {
	serializedData, err := so.Serialize()
	if err != nil {
		return err
	}
	err = sm.redisClient.Set(context.Background(), so.EventID, serializedData, 0).Err()
	return err
}

var nowFunc = time.Now

func (so *StateObject) LogTransition(from, to string, sm *StateMachine) {
	log := StateTransitionLog{
		EventID:   so.EventID,
		FromState: from,
		ToState:   to,
		Timestamp: nowFunc(),
	}

	if sm.DebugLogging {
		// Get caller info only if DebugLogging is enabled
		pc, file, line, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()

		file = filepath.Base(file)
		logStr := fmt.Sprintf("{timestamp: %s, event_id: %s, from_state: %s, to_state: %s, func: %s, file: %s, line: %d}\n",
			log.Timestamp, log.EventID, log.FromState, log.ToState, funcName, file, line)
		// Structured logging with debug info to stdout
		so.Logger.Debug("log_transition", zap.String("transition", logStr))
	} else {
		// Structured logging without debug info
		logStr := fmt.Sprintf("{timestamp: %s, event_id: %s, from_state: %s, to_state: %s}\n",
			log.Timestamp, log.EventID, log.FromState, log.ToState)
		so.Logger.Debug("log_transition", zap.String("transition", logStr))

	}
}

func (so *StateObject) MarkProcessed() error {
	so.State = "Processed"
	return so.CommitToDisk()
}

// DecodeDataToObject decodes the Data field into an object passed as an argument.
func (so *StateObject) DecodeDataToObject(obj interface{}) error {
	bytes, err := json.Marshal(so.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, obj)
}

// EncodeObjectToData encodes an object passed as an argument into the Data field.
func (so *StateObject) EncodeObjectToData(obj interface{}) error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}
	so.Data = data
	return nil
}
