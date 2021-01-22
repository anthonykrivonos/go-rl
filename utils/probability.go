package utils

// Uniform returns a uniform random probability given support size.
func Uniform(support int) float32 {
	return float32(1 / support)
}
