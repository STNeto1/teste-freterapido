package utils_test

import (
	"testing"

	"github.com/stneto1/teste-freterapido/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRange_Valid(t *testing.T) {
	t.Parallel()

	from := 1
	to := 10
	result := utils.Range(from, to)

	assert.Equal(t,
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		result)
}

func TestRange_Invalid(t *testing.T) {
	t.Parallel()

	from := 10
	to := 1
	result := utils.Range(from, to)

	assert.Equal(t, []int{}, result)
}

func TestRange_Invalid_SamePoints(t *testing.T) {
	t.Parallel()

	from := 10
	to := 10
	result := utils.Range(from, to)

	assert.Equal(t, []int{}, result)
}

func TestRangeWithStep_Valid(t *testing.T) {
	t.Parallel()

	from := 1
	to := 10
	step := 2

	assert.Equal(t, []int{1, 3, 5, 7, 9}, utils.RangeWithStep(from, to, step))
}
