package helpers

func Ptr[T any](in T) *T {
	return &in
}