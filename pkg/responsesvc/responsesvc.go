package responsesvc

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/expandorg/verification/pkg/verification"
)

type ResponseSVC interface {
	GetPending(jobID uint64, taskID uint64) (verification.TaskResponses, error)
}

type PendingResult struct {
	Responses verification.TaskResponses `json:"responses"`
}

func NormalizeRawMessage(raw json.RawMessage) (string, error) {
	var parsed map[string]interface{}
	err := json.Unmarshal(raw, &parsed)
	if err != nil {
		return "", err
	}

	for key, val := range parsed {
		if slice, ok := val.([]interface{}); ok {
			sort.Slice(slice, func(i, j int) bool { return toString(slice[i]) < toString(slice[j]) })
			parsed[key] = slice
		}
	}

	marshaled, err := json.Marshal(parsed)
	if err != nil {
		return "", err
	}

	return string(marshaled), nil
}

func toString(n interface{}) string {
	return fmt.Sprintf("%v", n)
}
