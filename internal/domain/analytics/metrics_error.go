package analytics

type AnalyticsInvalidLastQuotesError struct {
	Message string
}

func (e AnalyticsInvalidLastQuotesError) Error() string {
	return e.Message
}
