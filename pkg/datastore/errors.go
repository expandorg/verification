package datastore

type AlreadyResponded struct{}

func (err AlreadyResponded) Error() string {
	return "Response is already scored"
}

type AlreadyHasSettings struct{}

func (err AlreadyHasSettings) Error() string {
	return "Job already has settings"
}

type NoRowErr struct{}

func (err NoRowErr) Error() string {
	return "No Records"
}

type AlreadyAssigned struct{}

func (err AlreadyAssigned) Error() string {
	return "User is already assigned a task from this job"
}
