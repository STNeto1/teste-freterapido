package system

type SystemError struct {
	Message     string
	ShouldRetry bool
}

func (e SystemError) Error() string {
	return e.Message
}
