package automatic

type Consensus interface {
	Verify() (bool, error)
}

type consensus struct {
	store datastore.Storage
}

func NewConsensus(s datastore.Storage) Consensus {
	return &consenus{
		store:      s,
		authorizor: a,
		registry:   r,
		external:   e,
	}
}

func (s *consensus) Verify() (bool, error) {
	return true, nil
}
