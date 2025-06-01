package sl_utils

func Map[T any](s []T, callback func(el T, idx int, sliceLen int) T) []T {
	newSlice := make([]T, len(s), len(s))
	for idx, el := range s {
		newSlice[idx] = callback(el, idx, len(s))
	}

	return newSlice
}
