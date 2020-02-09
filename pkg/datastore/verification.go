package datastore

import (
	"strconv"
	"strings"

	"github.com/gemsorg/verification/pkg/verification"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	GetResponseAssignment(responseID uint64, verifierID uint64) (*verification.Assignment, error)
	Unassign(ID uint64) (*verification.Assignment, error)
	GetResponses(verification.Params) (verification.Responses, error)
	GetResponse(id string) (*verification.Response, error)
	CreateResponse(r verification.NewResponse) (*verification.Response, error)
	GetSettings(jobID uint64) (*verification.Settings, error)
	CreateSettings(s verification.Settings) (*verification.Settings, error)
	GetWhitelist(jobID uint64, verifierID uint64) (*verification.Whitelist, error)
	CreateAssignment(*verification.NewAssignment) (*verification.Assignment, error)
	GetAssignment(id string) (*verification.Assignment, error)
	GetAssignments(verification.Params) (verification.Assignments, error)
}

type VerificationStore struct {
	DB *sqlx.DB
}

func NewDatastore(db *sqlx.DB) *VerificationStore {
	return &VerificationStore{
		DB: db,
	}
}

func (vs *VerificationStore) GetResponses(p verification.Params) (verification.Responses, error) {
	responses := verification.Responses{}

	query := "SELECT * FROM verification_responses"
	conditions, args := p.ToQueryCondition()
	if len(args) > 0 {
		query = query + " WHERE " + conditions
	}
	err := vs.DB.Select(&responses, query, args...)
	if err != nil {
		return responses, err
	}
	return responses, nil
}

func (vs *VerificationStore) GetResponse(id string) (*verification.Response, error) {
	r := &verification.Response{}
	err := vs.DB.Get(r, "SELECT * FROM verification_responses WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (vs *VerificationStore) CreateResponse(r verification.NewResponse) (*verification.Response, error) {
	result, err := vs.DB.Exec(
		"INSERT INTO verification_responses (`job_id`, `task_id`, `response_id`, `worker_id`, `verifier_id`, `accepted`, `reason`) VALUES (?,?,?,?,?,?,?)",
		r.JobID, r.TaskID, r.ResponseID, r.WorkerID, r.VerifierID, r.Accepted, r.Reason,
	)

	if err != nil {
		mysqlerr, ok := err.(*mysql.MySQLError)
		// duplicate entry worker_id & job_id
		if ok && mysqlerr.Number == 1062 {
			return nil, AlreadyResponded{}
		}
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return vs.GetResponse(strconv.FormatInt(id, 10))
}

func (vs *VerificationStore) GetSettings(jobID uint64) (*verification.Settings, error) {
	set := []*verification.Settings{}
	err := vs.DB.Select(&set, "SELECT * FROM settings WHERE job_id = ?", jobID)

	if err != nil {
		return nil, err
	}
	if len(set) == 0 {
		return nil, NoRowErr{}
	}
	return set[0], nil
}

func (vs *VerificationStore) CreateSettings(s verification.Settings) (*verification.Settings, error) {
	// WE always replace the settings with the incoming
	_, err := vs.DB.Exec("REPLACE INTO settings (`job_id`, `manual`, `requester`, `limit`, `whitelist`, `agreement_count`) VALUES (?,?,?,?,?,?)", s.JobID, s.Manual, s.Requester, s.Limit, s.Whitelist, s.AgreementCount)
	if err != nil {
		mysqlerr, ok := err.(*mysql.MySQLError)
		// duplicate entry job_id
		if ok && mysqlerr.Number == 1062 {
			return nil, AlreadyHasSettings{}
		}
		return nil, err
	}
	return vs.GetSettings(s.JobID)
}

func (vs *VerificationStore) GetWhitelist(jobID uint64, verifierID uint64) (*verification.Whitelist, error) {
	wl := &verification.Whitelist{}
	err := vs.DB.Get(wl, "SELECT * FROM whitelists WHERE job_id = ? AND verifier_id = ?", jobID, verifierID)
	if err != nil {
		return nil, err
	}
	return wl, nil
}

func (vs *VerificationStore) GetAssignment(id string) (*verification.Assignment, error) {
	assignment := &verification.Assignment{}
	err := vs.DB.Get(assignment, "SELECT * FROM assignments WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	return assignment, nil
}

func (vs *VerificationStore) CreateAssignment(a *verification.NewAssignment) (*verification.Assignment, error) {
	result, err := vs.DB.Exec(
		"INSERT INTO assignments (job_id, task_id, verifier_id, active, expires_at) VALUES (?,?,?,?,DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 2 HOUR))",
		a.JobID, a.TaskID, a.VerifierID, 1)

	if err != nil {
		if err != nil {
			mysqlerr, ok := err.(*mysql.MySQLError)
			// duplicate entry verifier_id & job_id
			if ok && mysqlerr.Number == 1062 {
				return nil, AlreadyAssigned{}
			}
		}
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	assi, err := vs.GetAssignment(strconv.FormatInt(id, 10))

	if err != nil {
		return nil, err
	}

	return assi, nil
}

func (vs *VerificationStore) GetAssignments(p verification.Params) (verification.Assignments, error) {
	assignments := verification.Assignments{}
	query := "SELECT * FROM assignments"
	paramsQuery := []string{}
	args := []interface{}{}

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

	if len(paramsQuery) > 0 {
		query = query + " Where " + strings.Join(paramsQuery, " AND ")
	}

	err := vs.DB.Select(&assignments, query, args...)
	if err != nil {
		return assignments, err
	}
	return assignments, nil
}

func (vs *VerificationStore) GetResponseAssignment(responseID uint64, verifierID uint64) (*verification.Assignment, error) {
	return nil, nil
}

func (vs *VerificationStore) Unassign(ID uint64) (*verification.Assignment, error) {
	return nil, nil
}
