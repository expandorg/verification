package verification

type InvalidVerificationType struct {
	Manual bool
}

func (err InvalidVerificationType) Error() string {
	if err.Manual {
		return "This job has manual verifcaiton"
	}
	return "This job has automatic verifcaiton"
}
