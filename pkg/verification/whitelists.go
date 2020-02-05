package verification

// Whitelist for manual verification
type Whitelist struct {
	ID         uint64 `db:"id"`
	JobID      uint64 `db:"job_id"`
	VerifierID uint64 `db:"verifier_id"`
}
