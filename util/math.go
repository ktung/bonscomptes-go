package util

import "math"

const epsilon = 1e-2 // Tolerance for floating-point comparison

func IsZero(value float64) bool {
	return math.Abs(value) < epsilon
}
