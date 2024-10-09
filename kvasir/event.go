package kvasir

import (
	abci "github.com/cometbft/cometbft/abci/types"
)

type rawRequest struct {
	contract  string
	requestID uint64
	calldata  string
}

// GetEventValues returns the list of all values in the given log with the given type and key.
func GetEventValues(log []abci.Event, evType string, evKey string) (res []string) {
	for _, ev := range log {
		if ev.Type != evType {
			continue
		}

		for _, attr := range ev.Attributes {
			if attr.Key == evKey {
				res = append(res, attr.Value)
			}
		}
	}
	return res
}
