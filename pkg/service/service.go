package service

import (
	"github.com/gemsorg/verification/pkg/authentication"
	"github.com/gemsorg/verification/pkg/authorization"
	"github.com/gemsorg/verification/pkg/datastore"
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
}

func New(s datastore.Storage, a authorization.Authorizer) *service {
	return &service{
		store:      s,
		authorizor: a,
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

	resp, err := s.CreateResponse(r)
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

	return resp, nil
}
