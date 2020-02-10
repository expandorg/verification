package automatic

import (
	"github.com/gemsorg/verification/pkg/datastore"
	"github.com/gemsorg/verification/pkg/responsesvc"
	"github.com/gemsorg/verification/pkg/verification"
)

type Consensus interface {
	Verify(r verification.NewResponse, set *verification.Settings) (bool, error)
}

type ResponsesMap map[string]responsesvc.Responses

type consensus struct {
	store       datastore.Storage
	responseSVC responsesvc.ResponseSVC
}

func NewConsensus(s datastore.Storage, rs responsesvc.ResponseSVC) Consensus {
	return &consensus{
		store:       s,
		responseSVC: rs,
	}
}

func (s *consensus) Verify(r verification.NewResponse, set *verification.Settings) (bool, error) {
	responses, err := s.responseSVC.GetPending(r.JobID, r.TaskID)
	if err != nil {
		return false, err
	}
	if int64(len(responses)) < set.AgreementCount.Int64 {
		return false, err
	}
	responseMap, err := groupByRawMessage(responses)
	if err != nil {
		return false, err
	}

	consensus := responseMap.Conesensus(set.AgreementCount.Int64)

	if consensus == nil {
		return false, nil
	}
	return true, nil
}

func (rm ResponsesMap) Conesensus(agreementCount int64) responsesvc.Responses {
	var leader responsesvc.Responses = nil
	var leaderLen int64 = 0

	for _, responses := range rm {
		ln := int64(len(responses))
		if ln > leaderLen && ln >= agreementCount {
			leaderLen = ln
			leader = responses
		}
	}
	return leader
}

func groupByRawMessage(rs responsesvc.Responses) (ResponsesMap, error) {
	result := ResponsesMap{}
	for _, r := range rs {
		normalized, err := responsesvc.NormalizeRawMessage(r.Value)
		if err != nil {
			return result, err
		}

		if result[normalized] == nil {
			result[normalized] = responsesvc.Responses{r}
		} else {
			result[normalized] = append(result[normalized], r)
		}
	}
	return result, nil
}
