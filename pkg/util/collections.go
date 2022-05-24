package util

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func Filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func Some[T comparable](elems []T, f func(T) bool) bool {
	filtered := Filter(elems, f)
	count := len(filtered)
	return count > 0
}
