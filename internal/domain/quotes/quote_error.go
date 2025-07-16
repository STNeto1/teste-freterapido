package quotes

type QuoteRequestError struct {
	Message string
}

func (q QuoteRequestError) Error() string {
	return q.Message
}

type QuoteInvalidZipcodeError struct {
	Message string
}

func (q QuoteInvalidZipcodeError) Error() string {
	return q.Message
}

type QuoteRequestErrorSetError struct {
	Errors []string
}

// Error implements the error interface
// Not for public use (as a serializable error)
func (q QuoteRequestErrorSetError) Error() string {
	return "<SHOULD_NEVER_LEAVE_SYSTEM>"
}
