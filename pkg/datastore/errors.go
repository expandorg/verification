package datastore

import (
	"fmt"

	"github.com/expandorg/verification/pkg/nulls"
)

type AlreadyResponded struct{}

func (err AlreadyResponded) Error() string {
	return "Response is already scored"
}

type AlreadyHasSettings struct{}

func (err AlreadyHasSettings) Error() string {
	return "Job already has settings"
}

type NoRowErr struct{}

func (err NoRowErr) Error() string {
	return "No Records"
}

type AlreadyAssigned struct{}

func (err AlreadyAssigned) Error() string {
	return "User is already assigned a task from this job"
}

type NoAssignmentsAvailable struct {
	JobID uint64
}

func (err NoAssignmentsAvailable) Error() string {
	return fmt.Sprintf("There are no available verification for job: %d", err.JobID)
}

type AssignmentNotFound struct {
	ID         string
	WorkerID   uint64
	VerifierID nulls.Int64
	JobID      uint64
	ResponseID uint64
}

func (err AssignmentNotFound) Error() string {
	return fmt.Sprintf("No Record found for id: %s, worker_id: %d, verifier_id: %d, job_id: %d, response_id: %d", err.ID, err.WorkerID, err.VerifierID.Int64, err.JobID, err.ResponseID)
}
