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
	err := quotes.QuoteInvalidZipcode{
		Message: "Invalid zipcode",
	}

	assert.Equal(t, "Invalid zipcode", err.Error())
}
