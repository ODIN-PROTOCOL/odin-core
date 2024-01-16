package yoda

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

type rawRequest struct {
	dataSourceID   types.DataSourceID
	dataSourceHash string
	externalID     types.ExternalID
	calldata       string
}

// GetEventValues returns the list of all values in the given log with the given type and key.
func GetEventValues(log sdk.ABCIMessageLog, evType string, evKey string) (res []string) {
	for _, ev := range log.Events {
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

// GetEventValue checks and returns the exact value in the given log with the given type and key.
func GetEventValue(log sdk.ABCIMessageLog, evType string, evKey string) (string, error) {
	values := GetEventValues(log, evType, evKey)
	if len(values) == 0 {
		return "", fmt.Errorf("cannot find event with type: %s, key: %s", evType, evKey)
	}
	if len(values) > 1 {
		return "", fmt.Errorf("found more than one event with type: %s, key: %s", evType, evKey)
	}
	return values[0], nil
}
