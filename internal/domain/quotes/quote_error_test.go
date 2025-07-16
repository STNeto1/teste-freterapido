package quotes_test

import (
	"testing"

	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
	"github.com/stretchr/testify/assert"
)

func TestQuoteInvalidCategory(t *testing.T) {
	err := quotes.QuoteInvalidCategoryError{
		Message: "Invalid category",
	}

	assert.Equal(t, "Invalid category", err.Error())
}

func TestQuoteInvalidZipcode(t *testing.T) {
	err := quotes.QuoteInvalidZipcodeError{
		Message: "Invalid zipcode",
	}

	assert.Equal(t, "Invalid zipcode", err.Error())
}
func TestQuoteQuoteInvalidDimensionError(t *testing.T) {
	err := quotes.QuoteInvalidDimensionError{
		Message: "SUT",
	}

	assert.Equal(t, "SUT", err.Error())
}

func TestQuoteQuoteInvalidAmountError(t *testing.T) {
	err := quotes.QuoteInvalidAmountError{
		Message: "SUT",
	}

	assert.Equal(t, "SUT", err.Error())
}

func TestQuoteQuoteInvalidPriceError(t *testing.T) {
	err := quotes.QuoteInvalidPriceError{
		Message: "SUT",
	}

	assert.Equal(t, "SUT", err.Error())
}

func TestQuoteQuoteInvalidWeightError(t *testing.T) {
	err := quotes.QuoteInvalidWeightError{
		Message: "SUT",
	}

	assert.Equal(t, "SUT", err.Error())
}

func TestQuoteRequestError(t *testing.T) {
	err := quotes.QuoteRequestError{
		Message: "SUT",
	}

	assert.Equal(t, "SUT", err.Error())
}
