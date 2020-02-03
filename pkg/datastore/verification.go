package datastore

import (
	"strconv"

	"github.com/gemsorg/verification/pkg/verification"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	GetResponses(verification.Params) (verification.Responses, error)
	GetResponse(id string) (*verification.Response, error)
	CreateResponse(r verification.NewResponse) (*verification.Response, error)
	GetSettings(jobID uint64) (*verification.Settings, error)
	CreateSettings(s verification.Settings) (*verification.Settings, error)
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
	condition, args := p.ToQueryCondition()
	if len(args) > 0 {
		query = query + " WHERE " + condition
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
		"INSERT INTO verification_responses (`job_id`, `task_id`, `response_id`, `worker_id`, `verifier_id`, `value`, `reason`) VALUES (?,?,?,?,?,?,?)",
		r.JobID, r.TaskID, r.ResponseID, r.WorkerID, r.VerifierID, r.Value, r.Reason,
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
	_, err := vs.DB.Exec("REPLACE INTO settings (job_id) VALUES (?)", s.JobID)
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
