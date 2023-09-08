package main

import "encoding/json"

func Deserialize(data []byte) (*StateObject, error) {
	var so StateObject
	err := json.Unmarshal(data, &so)
	return &so, err
}

// Helper functions for creating chains
func CreateHandlerChain(handlers ...Handler) Handler {
	for i := 0; i < len(handlers)-1; i++ {
		handlers[i].SetNext(handlers[i+1])
	}
	return handlers[0]
}
