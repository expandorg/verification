package service

import (
	"github.com/expandorg/verification/pkg/authentication"
	"github.com/expandorg/verification/pkg/authorization"
	"github.com/expandorg/verification/pkg/automatic"
	"github.com/expandorg/verification/pkg/datastore"
	"github.com/expandorg/verification/pkg/externalsvc"
	"github.com/expandorg/verification/pkg/registrysvc"
	"github.com/expandorg/verification/pkg/verification"
)

type VerificationService interface {
	Healthy() bool
	SetAuthData(data authentication.AuthData)

	CreateEmptyAssignment(r verification.TaskResponse, set *verification.Settings) (*verification.Assignment, error)
	GetAssignments(verification.Params) (verification.Assignments, error)
	GetAssignment(id string) (*verification.Assignment, error)
	Assign(r verification.NewAssignment, set *verification.Settings) (*verification.Assignment, error)
	DeleteAssignment(id string) (bool, error)
	UpdateAssignment(a verification.Assignment) (*verification.Assignment, error)

	VerifyManual(r verification.NewVerificationResponse, set *verification.Settings) (*verification.VerificationResponse, error)
	VerifyAutomatic(r verification.TaskResponse, set *verification.Settings) (verification.VerificationResponses, error)

	GetResponses(verification.Params) (verification.VerificationResponses, error)
	GetResponse(id string) (*verification.VerificationResponse, error)

	GetSettings(jobID uint64) (*verification.Settings, error)
	CreateSettings(verification.Settings) (*verification.Settings, error)

	GetJobsWithEmptyAssignments() (verification.JobEmptyAssignments, error)
}

type service struct {
	store      datastore.Storage
	authorizor authorization.Authorizer
	registry   registrysvc.RegistrySVC
	external   externalsvc.External
	consensus  automatic.Consensus
}

func New(s datastore.Storage, a authorization.Authorizer, r registrysvc.RegistrySVC, e externalsvc.External, c automatic.Consensus) *service {
	return &service{
		store:      s,
		authorizor: a,
		registry:   r,
		external:   e,
		consensus:  c,
	}
}

func (s *service) Healthy() bool {
	return true
}

func (s *service) SetAuthData(data authentication.AuthData) {
	s.authorizor.SetAuthData(data)
}

func (s *service) GetResponses(p verification.Params) (verification.VerificationResponses, error) {
	return s.store.GetResponses(p)
}

func (s *service) GetResponse(id string) (*verification.VerificationResponse, error) {
	return s.store.GetResponse(id)
}

func (s *service) GetSettings(jobID uint64) (*verification.Settings, error) {
	set, err := s.store.GetSettings(jobID)
	if _, ok := err.(datastore.NoRowErr); ok {
		return nil, nil
	}
	return set, err
}

func (s *service) CreateSettings(set verification.Settings) (*verification.Settings, error) {
	return s.store.CreateSettings(set)
}

func (s *service) Assign(a verification.NewAssignment, set *verification.Settings) (*verification.Assignment, error) {
	asgn := verification.New(&a, nil)
	allowed, err := asgn.IsAllowed(set)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, AssignmentNotAllowed{}
	}

	return nil, nil
	// return s.store.CreateAssignment(&a)
}

func (s *service) CreateEmptyAssignment(r verification.TaskResponse, set *verification.Settings) (*verification.Assignment, error) {
	if set.Manual {
		return nil, InvalidVerificationType{set.Manual}
	}
	empty := verification.EmptyAssignment{
		ResponseID: r.ID,
		TaskID:     r.TaskID,
		JobID:      r.JobID,
	}
	return s.store.CreateAssignment(&empty)
}

func (s *service) DeleteAssignment(id string) (bool, error) {
	return s.store.DeleteAssignment(id)
}

func (s *service) UpdateAssignment(a verification.Assignment) (*verification.Assignment, error) {
	return s.store.UpdateAssignment(&a)
}

func (s *service) VerifyManual(r verification.NewVerificationResponse, set *verification.Settings) (*verification.VerificationResponse, error) {
	if !set.Manual {
		return nil, InvalidVerificationType{set.Manual}
	}

	// check that verifier is assigned
	assignment, err := s.store.GetAssignmentByResponseAndVerifier(r.ResponseID, r.VerifierID)
	if err != nil {
		return nil, Uniassigned{}
	}
	resp, err := s.store.CreateResponse(r.ToVerificationResponse())
	if err != nil {
		return nil, err
	}
	// unassign verification
	assignment.Status = verification.InActive
	_, err = s.store.UpdateAssignment(assignment)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *service) VerifyAutomatic(r verification.TaskResponse, set *verification.Settings) (verification.VerificationResponses, error) {
	if set.Manual {
		return nil, InvalidVerificationType{set.Manual}
	}

	vrs, err := s.callAutomaticVerification(r, set)
	if err != nil {
		return nil, err
	}

	if len(vrs) == 0 {
		// Nothing to verify
		return verification.VerificationResponses{}, nil
	}

	// create responses
	// TODO: implement response verification (jobId, taskId, responseID, workerId)
	responses, err := s.store.CreateResponses(vrs.ToVerificationResponses())
	if err != nil {
		return nil, err
	}
	return responses, nil
}

func (s *service) callAutomaticVerification(r verification.TaskResponse, set *verification.Settings) (externalsvc.VerificationResults, error) {
	reg := s.GetRegistration(r.JobID, registrysvc.ResponseVerifier)
	if reg != nil {
		return s.external.Verify(reg, r)
	}
	return s.consensus.Verify(r, set, s.authorizor.GetAuthToken())
}

func (s *service) GetRegistration(jobID uint64, svcType string) *registrysvc.Registration {
	r, _ := s.registry.GetRegistration(jobID)
	if r != nil && r.Services[svcType] != nil {
		return r
	}
	return nil
}

func (s *service) GetAssignments(p verification.Params) (verification.Assignments, error) {
	return s.store.GetAssignments(p)
}

func (s *service) GetAssignment(id string) (*verification.Assignment, error) {
	return s.store.GetAssignment(id)
}

func (s *service) GetJobsWithEmptyAssignments() (verification.JobEmptyAssignments, error) {
	return s.store.GetJobsWithEmptyAssignments()
}
