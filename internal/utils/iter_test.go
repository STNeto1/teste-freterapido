package utils_test

import (
	"testing"

	"github.com/stneto1/teste-freterapido/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4, 6, 8, 10}

	result := utils.Map(input, func(i int) int {
		return i * 2
	})

	assert.Equal(t, expected, result)
}
