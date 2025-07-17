package system_test

import (
	"testing"

	"github.com/stneto1/teste-freterapido/internal/domain/system"
	"github.com/stretchr/testify/assert"
)

func TestSystemError_Error(t *testing.T) {
	err := system.SystemError{
		Message: "SUT",
	}

	assert.Equal(t, "SUT", err.Error())
}
