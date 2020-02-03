package verification

type Settings struct {
	ID    uint64 `json:"id" db:"id"`
	JobID uint64 `json:"job_id" db:"job_id"`
}
