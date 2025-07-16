package utils

// Range creates a slice of integers from a to b
func Range(a, b int) []int {
	return RangeWithStep(a, b, 1)
}

func RangeWithStep(a, b, step int) []int {
	if b < a {
		return []int{}
	}

	if a == b {
		return []int{}
	}

	size := (b-a)/step + 1
	nums := make([]int, size)
	for i := range size {
		nums[i] = a + i*step
	}
	return nums
}
