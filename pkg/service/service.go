package service

import (
	"github.com/gemsorg/verification/pkg/authorization"
	"github.com/gemsorg/verification/pkg/datastore"
)

type VerificationService interface {
	Healthy() bool
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
