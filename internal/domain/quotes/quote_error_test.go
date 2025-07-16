package quotes_test

import (
	"testing"

	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
	"github.com/stretchr/testify/assert"
)

func TestQuoteRequestError(t *testing.T) {
	err := quotes.QuoteRequestError{
		Message: "SUT",
	}

	assert.Equal(t, "SUT", err.Error())
}

func TestQuoteRequestErrorSetError(t *testing.T) {
	err := quotes.QuoteRequestErrorSetError{
		Errors: []string{},
	}

	assert.Equal(t, "<SHOULD_NEVER_LEAVE_SYSTEM>", err.Error())
}
