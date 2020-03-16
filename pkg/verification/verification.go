package verification

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/expandorg/verification/pkg/nulls"
)

type Status string

const (
	Active   Status = "active"
	InActive Status = "inactive"
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Rejected Status = "rejected"
	Expired  Status = "expired"
)

// VerificationResponse onse entity
type VerificationResponse struct {
	ID         uint64       `db:"id" json:"id"`
	JobID      uint64       `db:"job_id" json:"job_id"`
	TaskID     uint64       `db:"task_id" json:"task_id"`
	ResponseID uint64       `db:"response_id" json:"response_id"`
	WorkerID   uint64       `db:"worker_id" json:"worker_id"`
	VerifierID nulls.Int64  `db:"verifier_id" json:"verifier_id"`
	Accepted   bool         `db:"accepted" json:"accepted"`
	Reason     nulls.String `db:"reason" json:"reason"`
	CreatedAt  time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at" json:"updated_at"`
}

// VerificationResponses list
type VerificationResponses []VerificationResponse

// NewVerificationResponse entity (params)
type NewVerificationResponse struct {
	JobID      uint64       `json:"job_id"`
	TaskID     uint64       `json:"task_id"`
	ResponseID uint64       `json:"response_id"`
	WorkerID   uint64       `json:"worker_id"`
	VerifierID int64        `json:"verifier_id"`
	Accepted   bool         `json:"accepted"`
	Reason     nulls.String `json:"reason"`
}

// TaskResponse entity (params)
type TaskResponse struct {
	ID        uint64          `json:"id"`
	JobID     uint64          `json:"job_id"`
	TaskID    uint64          `json:"task_id"`
	WorkerID  uint64          `json:"worker_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Value     json.RawMessage `json:"value"`
}

type TaskResponses []TaskResponse

// Params for querying responses
type Params struct {
	WorkerID   string
	VerifierID string
	JobID      string
	TaskID     string
	ResponseID string
	Status     Status
}

// ToQueryCondition converts Params to sql query condition
func (p Params) ToQueryCondition() (string, []interface{}) {
	paramsQuery := []string{}
	args := []interface{}{}
	if p.WorkerID != "" && p.WorkerID != "0" {
		args = append(args, p.WorkerID)
		paramsQuery = append(paramsQuery, "worker_id=?")
	}
	if p.VerifierID != "" && p.VerifierID != "0" {
		args = append(args, p.VerifierID)
		paramsQuery = append(paramsQuery, "verifier_id=?")
	}
	if p.JobID != "" && p.JobID != "0" {
		args = append(args, p.JobID)
		paramsQuery = append(paramsQuery, "job_id=?")
	}
	if p.TaskID != "" && p.TaskID != "0" {
		args = append(args, p.TaskID)
		paramsQuery = append(paramsQuery, "task_id=?")
	}
	if p.ResponseID != "" && p.ResponseID != "0" {
		args = append(args, p.ResponseID)
		paramsQuery = append(paramsQuery, "response_id=?")
	}
	if p.Status != "" && p.Status != "0" {
		args = append(args, p.Status)
		paramsQuery = append(paramsQuery, "status=?")
	}

	return strings.Join(paramsQuery, " AND "), args
}

func (n NewVerificationResponse) ToVerificationResponse() VerificationResponse {
	return VerificationResponse{
		JobID:      n.JobID,
		TaskID:     n.TaskID,
		ResponseID: n.ResponseID,
		WorkerID:   n.WorkerID,
		VerifierID: nulls.NewInt64(n.VerifierID),
		Accepted:   n.Accepted,
		Reason:     n.Reason,
	}
}

func (rs TaskResponses) Has(resp TaskResponse) bool {
	for _, r := range rs {
		if r.ID == resp.ID {
			return true
		}
	}
	return false
}
