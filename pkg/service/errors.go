package service

type AssignmentNotAllowed struct {
	Manual bool
}

func (err AssignmentNotAllowed) Error() string {
	return "This assignment is not allowed"
}

type InvalidVerificationType struct {
	Manual bool
}

func (err InvalidVerificationType) Error() string {
	if err.Manual {
		return "This job has manual verifcaiton"
	}
	return "This job has automatic verifcaiton"
}

type Uniassigned struct {
	Manual bool
}

func (err Uniassigned) Error() string {
	return "The user is not assigned to verificaton"
}
