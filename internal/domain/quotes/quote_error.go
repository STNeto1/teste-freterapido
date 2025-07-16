package quotes

type QuoteInvalidCategoryError struct {
	Message string
}

func (q QuoteInvalidCategoryError) Error() string {
	return q.Message
}

type QuoteInvalidZipcode struct {
	Message string
}

func (q QuoteInvalidZipcode) Error() string {
	return q.Message
}
