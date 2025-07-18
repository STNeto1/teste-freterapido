package system_test

import (
	"testing"
	"time"

	"github.com/stneto1/teste-freterapido/internal/domain/system"
	"github.com/stretchr/testify/assert"
)

func TestProcessQuoteServiceConfig(t *testing.T) {
	start := system.Start{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 1,
		TryQuotesRetries:  1,
		TryQuotesTimeout:  time.Millisecond,
		AddQuotesRetries:  1,
		AddQuotesTimeout:  time.Millisecond,
	}

	result := start.ProcessQuoteServiceConfig(nil)

	assert.Equal(t, "SUT", result.RegisteredNumber)
	assert.Equal(t, "SUT", result.Token)
	assert.Equal(t, "SUT", result.PlatformCode)
	assert.Equal(t, int64(1), result.DispatcherZipCode)
	assert.Equal(t, 1, result.TryQuotesRetries)
	assert.Equal(t, time.Millisecond, result.TryQuotesTimeout)
	assert.Equal(t, 1, result.AddQuotesRetries)
	assert.Equal(t, time.Millisecond, result.AddQuotesTimeout)
	assert.Nil(t, result.Logger)
}
