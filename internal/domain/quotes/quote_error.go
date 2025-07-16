package quotes

type QuoteInvalidCategoryError struct {
	Message string
}

func (q QuoteInvalidCategoryError) Error() string {
	return q.Message
}

type QuoteInvalidZipcodeError struct {
	Message string
}

func (q QuoteInvalidZipcodeError) Error() string {
	return q.Message
}

type QuoteInvalidDimensionError struct {
	Message string
}

func (q QuoteInvalidDimensionError) Error() string {
	return q.Message
}

type QuoteInvalidAmountError struct {
	Message string
}

func (q QuoteInvalidAmountError) Error() string {
	return q.Message
}

type QuoteInvalidPriceError struct {
	Message string
}

func (q QuoteInvalidPriceError) Error() string {
	return q.Message
}

type QuoteInvalidWeightError struct {
	Message string
}

func (q QuoteInvalidWeightError) Error() string {
	return q.Message
}
