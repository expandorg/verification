package responsesvc

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/gemsorg/verification/pkg/nulls"
)

type ResponseSVC interface {
	GetPending(jobID uint64, taskID uint64) (Responses, error)
}

type PendingResult struct {
	Responses Responses `json:"responses"`
}

type Response struct {
	ID         uint64          `json:"id"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	WorkerID   uint64          `json:"worker_id"`
	JobID      uint64          `json:"job_id"`
	TaskID     uint64          `json:"task_id"`
	Value      json.RawMessage `json:"value"`
	IsAccepted nulls.Bool      `json:"is_accepted"`
}

type Responses []Response

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
