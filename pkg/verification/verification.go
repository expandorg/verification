package verification

import (
	"strings"
	"time"

	"github.com/gemsorg/backend/nulls"
)

type Response struct {
	ID         uint64      `db:"id" json:"id"`
	JobID      uint64      `db:"job_id" json:"job_id"`
	TaskID     uint64      `db:"task_id" json:"task_id"`
	ResponseID uint64      `db:"response_id" json:"response_id"`
	WorkerID   uint64      `db:"worker_id" json:"worker_id"`
	VerifierID nulls.Int64 `db:"verifier_id" json:"verifier_id"`
	Value      bool        `db:"value" json:"value"`
	Reason     string      `db:"reason" json:"reason"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at" json:"updated_at"`
}

type Responses []Response

type NewResponse struct {
	JobID      uint64      `json:"job_id"`
	TaskID     uint64      `json:"task_id"`
	ResponseID uint64      `json:"response_id"`
	WorkerID   uint64      `json:"worker_id"`
	VerifierID nulls.Int64 `json:"verifier_id"`
	Value      bool        `json:"value"`
	Reason     string      `json:"reason"`
}

type Params struct {
	WorkerID string
	JobID    string
	TaskID   string
}

func (p Params) ToQueryCondition() (string, []interface{}) {
	paramsQuery := []string{}
	args := []interface{}{}
	if p.WorkerID != "" && p.WorkerID != "0" {
		args = append(args, p.WorkerID)
		paramsQuery = append(paramsQuery, "worker_id=?")
	}
	if p.JobID != "" && p.JobID != "0" {
		args = append(args, p.JobID)
		paramsQuery = append(paramsQuery, "job_id=?")
	}
	if p.TaskID != "" && p.TaskID != "0" {
		args = append(args, p.TaskID)
		paramsQuery = append(paramsQuery, "task_id=?")
	}

	return strings.Join(paramsQuery, " AND "), args
}
