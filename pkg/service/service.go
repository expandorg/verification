package service

import (
	"github.com/gemsorg/verification/pkg/authentication"
	"github.com/gemsorg/verification/pkg/authorization"
	"github.com/gemsorg/verification/pkg/datastore"
	"github.com/gemsorg/verification/pkg/externalsvc"
	"github.com/gemsorg/verification/pkg/registrysvc"
	"github.com/gemsorg/verification/pkg/verification"
)

type VerificationService interface {
	Healthy() bool
	SetAuthData(data authentication.AuthData)

	Assign(r verification.NewAssignment, set *verification.Settings) (*verification.Assignment, error)

	VerifyManual(r verification.NewResponse, set *verification.Settings) (*verification.Response, error)
	VerifyAutomatic(r verification.NewResponse, set *verification.Settings) (*verification.Response, error)

	GetResponses(verification.Params) (verification.Responses, error)
	GetResponse(id string) (*verification.Response, error)
	CreateResponse(n verification.NewResponse) (*verification.Response, error)

	GetSettings(jobID uint64) (*verification.Settings, error)
	CreateSettings(verification.Settings) (*verification.Settings, error)
}

type service struct {
	store      datastore.Storage
	authorizor authorization.Authorizer
	registry   registrysvc.RegistrySVC
	external   externalsvc.External
}

func New(s datastore.Storage, a authorization.Authorizer, r registrysvc.RegistrySVC, e externalsvc.External) *service {
	return &service{
		store:      s,
		authorizor: a,
		registry:   r,
		external:   e,
	}
}

func (s *service) Healthy() bool {
	return true
}

func (s *service) SetAuthData(data authentication.AuthData) {
	s.authorizor.SetAuthData(data)
}

func (s *service) GetResponses(p verification.Params) (verification.Responses, error) {
	return s.store.GetResponses(p)
}

func (s *service) GetResponse(id string) (*verification.Response, error) {
	return s.store.GetResponse(id)
}

func (s *service) CreateResponse(n verification.NewResponse) (*verification.Response, error) {
	return s.store.CreateResponse(n)
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

	return s.store.CreateAssignment(&a)
}

func (s *service) VerifyManual(r verification.NewResponse, set *verification.Settings) (*verification.Response, error) {
	if !set.Manual {
		return nil, InvalidVerificationType{set.Manual}
	}

	// TOOD: wrap to tx-scope
	assignment, err := s.store.GetResponseAssignment(r.ResponseID, uint64(r.VerifierID.Int64))
	if err != nil {
		return nil, Uniassigned{}
	}
	resp, err := s.CreateResponse(r)
	if err != nil {
		return nil, err
	}
	_, err = s.store.Unassign(assignment.ID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *service) VerifyAutomatic(r verification.NewResponse, set *verification.Settings) (*verification.Response, error) {
	if set.Manual {
		return nil, InvalidVerificationType{set.Manual}
	}
	resp, err := s.CreateResponse(r)
	if err != nil {
		return nil, err
	}

	// reg := s.GetRegistration(r.JobID, registrysvc.ResponseVerifier)
	// if reg != nil {
	// 	s.external.Verify()
	// } else {

	// }

	return resp, nil
}

func (s *service) GetRegistration(jobID uint64, svcType string) *registrysvc.Registration {
	r, _ := s.registry.GetRegistration(jobID)
	if r != nil && r.Services[svcType] != nil {
		return r
	}
	return nil
}
