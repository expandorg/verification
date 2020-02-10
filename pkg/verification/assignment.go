package verification

import (
	"time"

	"github.com/gemsorg/verification/pkg/nulls"
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
	VerifierID uint64      `db:"verifier_id" json:"verifier_id"`
	ResponseID nulls.Int64 `db:"response_id" json:"response_id"`
	Active     nulls.Bool  `db:"active" json:"active"`
	AssignedAt time.Time   `db:"assigned_at" json:"assigned_at"`
	ExpiresAt  nulls.Time  `db:"expires_at" json:"expires_at"`
}

type NewAssignment struct {
	JobID      uint64 `json:"job_id"`
	TaskID     uint64 `json:"task_id"`
	VerifierID uint64 `json:"verifier_id"`
}

type Assignments []Assignment

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
