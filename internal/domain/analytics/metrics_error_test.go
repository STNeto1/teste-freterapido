package analytics_test

import (
	"testing"

	"github.com/stneto1/teste-freterapido/internal/domain/analytics"
	"github.com/stretchr/testify/assert"
)

func TestAnalyticsInvalidLastQuotesError_Error(t *testing.T) {
	err := analytics.AnalyticsInvalidLastQuotesError{
		Message: "SUT",
	}

	assert.Equal(t, "SUT", err.Error())
}
