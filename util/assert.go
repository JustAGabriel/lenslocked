package util

// Must asserts that a possible error return type is nil
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
