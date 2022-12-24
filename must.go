package rr

func Must[T any](t T, err error) T {
	Check(err)
	return t
}
