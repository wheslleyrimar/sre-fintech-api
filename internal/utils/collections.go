package utils

// Map applies f to each element of slice and returns a new slice.
func Map[T, K any](input []T, f func(T) K) ([]K, error) {
	if input == nil {
		return nil, nil
	}
	out := make([]K, len(input))
	for i := range input {
		out[i] = f(input[i])
	}
	return out, nil
}

// Contains returns true if any element satisfies predicate.
func Contains[T any](input []T, predicate func(T) bool) bool {
	for _, elem := range input {
		if predicate(elem) {
			return true
		}
	}
	return false
}
