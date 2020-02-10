package automatic

import (
	"github.com/gemsorg/verification/pkg/datastore"
	"github.com/gemsorg/verification/pkg/verification"
)

type Consensus interface {
	Verify(r verification.NewResponse, set *verification.Settings) (bool, error)
}

type consensus struct {
	store datastore.Storage
}

func NewConsensus(s datastore.Storage) Consensus {
	return &consensus{
		store: s,
	}
}

func (s *consensus) Verify(r verification.NewResponse, set *verification.Settings) (bool, error) {
	return true, nil
}
