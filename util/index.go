package util

func NextIndex[T any](slc []T, cur int) int {
	return (cur + 1) % len(slc)
}

func PrevIndex[T any](slc []T, cur int) int {
	l := len(slc)
	return (cur - 1 + l) % l
}
