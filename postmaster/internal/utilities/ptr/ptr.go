package ptr

func Ptr[T any](v T) *T {
	return &v
}

func ValueOrDefault[T any](v *T, def T) T {
	if v == nil {
		return def
	}

	return *v
}

func ValueOrZero[T any](v *T) T {
	var zero T

	return ValueOrDefault(v, zero)
}
