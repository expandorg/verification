package automatic

import (
	"github.com/gemsorg/verification/pkg/datastore"
	"github.com/gemsorg/verification/pkg/externalsvc"
	"github.com/gemsorg/verification/pkg/responsesvc"
	"github.com/gemsorg/verification/pkg/verification"
)

type Consensus interface {
	Verify(r verification.TaskResponse, set *verification.Settings) (*externalsvc.VerifyResponse, error)
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

func (s *consensus) Verify(r verification.TaskResponse, set *verification.Settings) (*externalsvc.VerifyResponse, error) {
	results := externalsvc.VerifyResponse{}

	responses, err := s.responseSVC.GetPending(r.JobID, r.TaskID)
	if err != nil {
		return nil, err
	}
	if int64(len(responses)) < set.AgreementCount.Int64 {
		return nil, err
	}

	grouped, err := groupByRawMessage(responses)
	if err != nil {
		return nil, err
	}

	consensus := grouped.Conesensus(set.AgreementCount.Int64)
	if consensus == nil {
		return &results, nil
	}

	for _, rsp := range responses {
		r := externalsvc.VerificationResult{
			JobID:      rsp.JobID,
			TaskID:     rsp.TaskID,
			ResponseID: rsp.ID,
			Accepted:   consensus.Has(rsp),
		}
		results.Results = append(results.Results, r)
	}
	return &results, nil
}

func (rm ResponsesMap) Conesensus(agreementCount int64) responsesvc.Responses {
	var leaders responsesvc.Responses = nil
	var leadersLen int64 = 0

	for _, responses := range rm {
		ln := int64(len(responses))

		if ln > leadersLen && ln >= agreementCount {
			leadersLen = ln
			leaders = responses
		}
	}
	return leaders
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
