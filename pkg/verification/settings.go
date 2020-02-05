package verification

import (
	"github.com/gemsorg/verification/pkg/nulls"
)

// Settings for job verification
type Settings struct {
	ID             uint64      `json:"id" db:"id"`
	JobID          uint64      `json:"job_id" db:"job_id"`
	Manual         bool        `json:"manual" db:"manual"`                   // verification can be automatic or manual
	Requester      bool        `json:"requester" db:"requester"`             // manual (assign): requester only
	Limit          int         `json:"limit" db:"limit"`                     // manual (assign): total assignment limit per job per worker
	Whitelist      bool        `json:"whitelist" db:"whitelist"`             // manual (assign): a whitelist for workers
	AgreementCount nulls.Int64 `json:"agreement_count" db:"agreement_count"` // automatic (verify): agreeement count
}
