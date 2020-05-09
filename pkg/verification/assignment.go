package verification

import (
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

type Asgn interface {
	IsAllowed(set *Settings) (bool, error)
	GetAssignment() *Assignment
	GetNewAssignment() *NewAssignment
}

type assignment struct {
	NewAssignment *NewAssignment
	Assignment    *Assignment
}

type Assignment struct {
	ID         uint64      `db:"id" json:"id"`
	JobID      uint64      `db:"job_id" json:"job_id"`
	TaskID     uint64      `db:"task_id" json:"task_id"`
	WorkerID   uint64      `db:"worker_id" json:"worker_id"`
	VerifierID nulls.Int64 `db:"verifier_id" json:"verifier_id"`
	ResponseID nulls.Int64 `db:"response_id" json:"response_id"`
	Active     nulls.Bool  `db:"active" json:"active"`
	Status     Status      `db:"status" json:"status"`
	AssignedAt time.Time   `db:"assigned_at" json:"assigned_at"`
	ExpiresAt  nulls.Time  `db:"expires_at" json:"expires_at"`
}

type Assignments []Assignment

type NewAssignment struct {
	JobID      uint64 `json:"job_id"`
	TaskID     uint64 `json:"task_id"`
	VerifierID uint64 `json:"verifier_id"`
}

type EmptyAssignment struct {
	JobID      uint64 `json:"job_id"`
	WorkerID   uint64 `json:"worker_id"`
	TaskID     uint64 `json:"task_id"`
	ResponseID uint64 `json:"response_id"`
}

type JobEmptyAssignment struct {
	JobID          uint64 `db:"job_id" json:"job_id"`
	AvailableCount uint64 `db:"empty_count" json:"empty_count"`
}

type JobEmptyAssignments []JobEmptyAssignment

func New(na *NewAssignment, a *Assignment) *assignment {
	return &assignment{
		NewAssignment: na,
		Assignment:    a,
	}
}
func (a *assignment) GetNewAssignment() *NewAssignment {
	return a.NewAssignment
}

func (a *assignment) GetAssignment() *Assignment {
	return a.Assignment
}

func (a *assignment) IsAllowed(set *Settings) (bool, error) {
	if set != nil && !set.Manual {
		return false, InvalidVerificationType{set.Manual}
	}

	return true, nil
}

func (jes JobEmptyAssignments) JobIDs() []uint64 {
	ids := make([]uint64, len(jes))
	for i, t := range jes {
		ids[i] = t.JobID
	}
	return ids
}

func (js JobEmptyAssignments) Filter(cond func(JobEmptyAssignment) bool) JobEmptyAssignments {
	result := make(JobEmptyAssignments, 0)
	for _, tx := range js {
		if cond(tx) {
			result = append(result, tx)
		}
	}
	return result
}
