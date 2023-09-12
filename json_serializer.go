package statemachine

import (
	"encoding/json"
)

type JSONSerialization struct{}

func (j JSONSerialization) Serialize(so *StateObject) ([]byte, error) {
	return json.Marshal(so)
}

func (j JSONSerialization) Deserialize(data []byte) (*StateObject, error) {
	var so StateObject
	err := json.Unmarshal(data, &so)
	return &so, err
}
