package datastore

import (
	"strconv"

	"github.com/expandorg/verification/pkg/verification"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	GetAssignmentByResponseAndVerifier(responseID uint64, verifierID int64) (*verification.Assignment, error)
	UpdateAssignment(a *verification.Assignment) (*verification.Assignment, error)
	GetResponses(verification.Params) (verification.VerificationResponses, error)
	GetResponse(id string) (*verification.VerificationResponse, error)
	CreateResponse(r verification.VerificationResponse) (*verification.VerificationResponse, error)
	CreateResponses(rs verification.VerificationResponses) (verification.VerificationResponses, error)
	GetSettings(jobID uint64) (*verification.Settings, error)
	CreateSettings(s verification.Settings) (*verification.Settings, error)
	GetWhitelist(jobID uint64, verifierID uint64) (*verification.Whitelist, error)
	CreateAssignment(*verification.EmptyAssignment) (*verification.Assignment, error)
	GetAssignment(id string) (*verification.Assignment, error)
	GetAssignments(verification.Params) (verification.Assignments, error)
	DeleteAssignment(id string) (bool, error)
	GetJobsWithEmptyAssignments(verifierID uint64) (verification.JobEmptyAssignments, error)
	GetEligibleJobIDs(verifierID uint64, JobIDs []uint64) ([]uint64, error)
	Assign(a *verification.NewAssignment) (*verification.Assignment, error)
}

type VerificationStore struct {
	DB *sqlx.DB
}

func NewDatastore(db *sqlx.DB) *VerificationStore {
	return &VerificationStore{
		DB: db,
	}
}

func (vs *VerificationStore) GetResponses(p verification.Params) (verification.VerificationResponses, error) {
	responses := verification.VerificationResponses{}

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

func (vs *VerificationStore) GetResponse(id string) (*verification.VerificationResponse, error) {
	return vs.getResponse(vs.DB, id)
}

func (vs *VerificationStore) getResponse(db DbQueryExecutor, id string) (*verification.VerificationResponse, error) {
	r := &verification.VerificationResponse{}
	err := db.Get(r, "SELECT * FROM verification_responses WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (vs *VerificationStore) CreateResponse(r verification.VerificationResponse) (*verification.VerificationResponse, error) {
	return vs.createResponse(vs.DB, r)
}

func (vs *VerificationStore) CreateResponses(rs verification.VerificationResponses) (verification.VerificationResponses, error) {
	responses := make(verification.VerificationResponses, len(rs))

	tx, err := vs.DB.Beginx()
	if err != nil {
		return nil, err
	}
	for i, r := range rs {
		vr, err := vs.createResponse(tx, r)
		if err != nil {
			return nil, vs.tryRollback(tx, err)
		}
		responses[i] = *vr
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return responses, nil
}

func (vs *VerificationStore) createResponse(db DbQueryExecutor, r verification.VerificationResponse) (*verification.VerificationResponse, error) {
	result, err := db.Exec(
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
	return vs.getResponse(db, strconv.FormatInt(id, 10))
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

func (vs *VerificationStore) GetAssignmentByResponseAndVerifier(responseID uint64, verifierID int64) (*verification.Assignment, error) {
	assignment := &verification.Assignment{}
	err := vs.DB.Get(assignment, "SELECT * FROM assignments WHERE response_id = ? AND verifier_id = ?", responseID, verifierID)
	if err != nil {
		return nil, err
	}
	return assignment, nil
}

func (vs *VerificationStore) CreateAssignment(a *verification.EmptyAssignment) (*verification.Assignment, error) {
	result, err := vs.DB.Exec(
		"INSERT INTO assignments (job_id, task_id, response_id, worker_id, active) VALUES (?,?,?,?,?)",
		a.JobID, a.TaskID, a.ResponseID, a.WorkerID, 0,
	)
	if err != nil {
		mysqlerr, ok := err.(*mysql.MySQLError)
		// duplicate entry verifier_id & job_id
		if ok && mysqlerr.Number == 1062 {
			return nil, AlreadyAssigned{}
		}
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return vs.GetAssignment(strconv.FormatInt(id, 10))
}

func (vs *VerificationStore) Assign(a *verification.NewAssignment) (*verification.Assignment, error) {
	result, err := vs.DB.Exec(
		"UPDATE assignments SET id=(SELECT @updated_id := id), verifier_id=?, active=?, status=?, assigned_at=CURRENT_TIMESTAMP, expires_at=DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 2 HOUR) WHERE job_id = ? AND verifier_id IS NULL AND worker_id <> ? LIMIT 1",
		a.VerifierID, 1, verification.Active, a.JobID, a.VerifierID,
	)
	if err != nil {
		mysqlerr, ok := err.(*mysql.MySQLError)
		// duplicate entry verifier_id & job_id
		if ok && mysqlerr.Number == 1062 {
			return nil, AlreadyAssigned{}
		}
		return nil, err
	}

	num, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, NoAssignmentsAvailable{a.JobID}
	}

	var assignmentID string
	err = vs.DB.Get(&assignmentID, "SELECT @updated_id")
	if err != nil {
		return nil, err
	}
	return vs.GetAssignment(assignmentID)
}

func (vs *VerificationStore) GetAssignments(p verification.Params) (verification.Assignments, error) {
	assignments := verification.Assignments{}
	query := "SELECT * FROM assignments"
	args := []interface{}{}

	conditions, args := p.ToQueryCondition()
	if len(args) > 0 {
		query = query + " WHERE " + conditions
	}

	err := vs.DB.Select(&assignments, query, args...)
	if err != nil {
		return assignments, err
	}

	return assignments, nil
}

func (vs *VerificationStore) DeleteAssignment(id string) (bool, error) {
	result, err := vs.DB.Exec("DELETE FROM assignments WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	numAffected, err := result.RowsAffected()

	if numAffected == 0 {
		return false, AssignmentNotFound{ID: id}
	}

	return true, nil
}

func (vs *VerificationStore) UpdateAssignment(a *verification.Assignment) (*verification.Assignment, error) {
	if a.JobID == 0 || !a.VerifierID.Valid || !a.ResponseID.Valid || a.Status == "" {
		return nil, AssignmentNotFound{VerifierID: a.VerifierID, JobID: a.JobID, ResponseID: uint64(a.ResponseID.Int64)}
	}

	var active bool
	if a.Status == verification.Active {
		active = true
	}

	query := "UPDATE assignments SET status=?, active=? WHERE verifier_id=? AND job_id=? AND response_id=?"
	_, err := vs.DB.Exec(query, a.Status, active, a.VerifierID, a.JobID, a.ResponseID)
	if err != nil {
		return nil, err
	}

	if a.ID == 0 {
		return a, nil
	}
	return vs.GetAssignment(strconv.FormatUint(a.ID, 10))
}

func (vs *VerificationStore) GetJobsWithEmptyAssignments(verifierID uint64) (verification.JobEmptyAssignments, error) {
	a := verification.JobEmptyAssignments{}
	err := vs.DB.Select(&a, "SELECT job_id, count(id) as empty_count FROM assignments WHERE verifier_id is NULL and worker_id <> ? GROUP BY job_id ", verifierID)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (vs *VerificationStore) GetEligibleJobIDs(verifierID uint64, JobIDs []uint64) ([]uint64, error) {
	query, args, err := sqlx.In(`
		SELECT t.job_id as job_id from 
		(
			SELECT s.job_id from settings s 
			JOIN whitelists w on s.job_id = w.job_id and w.verifier_id = ?
			WHERE s.manual = TRUE and s.requester = FALSE and s.whitelist = TRUE			
				UNION			
			SELECT s.job_id from settings s 
			WHERE s.manual = TRUE and s.requester = FALSE and s.whitelist = FALSE
		) as t 
		WHERE t.job_id in (?)	
		`,
		verifierID,
		JobIDs,
	)
	eligible := []uint64{}
	query = vs.DB.Rebind(query)
	err = vs.DB.Select(&eligible, query, args...)
	if err != nil {
		return nil, err
	}
	return eligible, nil
}

type DbQueryExecutor interface {
	sqlx.Execer
	sqlx.Queryer
	Get(dest interface{}, query string, args ...interface{}) error
}

func (vs *VerificationStore) tryRollback(tx *sqlx.Tx, err error) error {
	roolbackErr := tx.Rollback()
	if roolbackErr != nil {
		return roolbackErr
	}
	return err
}
