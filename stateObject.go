package statemachine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type StateObject struct {
	Data       map[string]interface{}                           `json:"data"`
	State      string                                           `json:"state"`
	EventID    string                                           `json:"eventID"`
	Logger     func(format string, a ...any) (n int, err error) `json:"-"`
	CommitFunc func() error                                     `json:"-"`
}

func NewStateObject(data map[string]interface{}, sm *StateMachine) *StateObject {
	var state = &StateObject{
		Data:   data,
		State:  SIMNotActivated, // or some other default state
		Logger: fmt.Printf,
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
		// Structured logging with debug info to stdout
		so.Logger(
			"{timestamp: %s, event_id: %s, from_state: %s, to_state: %s, func: %s, file: %s, line: %d}\n",
			log.Timestamp, log.EventID, log.FromState, log.ToState, funcName, file, line,
		)
	} else {
		// Structured logging without debug info
		so.Logger(
			"{timestamp: %s, event_id: %s, from_state: %s, to_state: %s}\n",
			log.Timestamp, log.EventID, log.FromState, log.ToState,
		)
	}
}

func (so *StateObject) MarkProcessed() error {
	so.State = "Processed"
	return so.CommitToDisk()
}
