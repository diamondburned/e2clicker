package ptr

// To returns a pointer to the value passed as argument.
func To[T any](v T) *T {
	return &v
}

// Deref returns the value that v points to, or the zero value of T if v is nil.
func Deref[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// DerefOr returns the value that v points to, or the value passed as or if v is
// nil.
func DerefOr[T any](v *T, or T) T {
	if v == nil {
		return or
	}
	return *v
}

// ToIf returns a pointer to the value passed as argument if ok is true, or nil
// if ok is false.
func ToIf[T any](v T, ok bool) *T {
	if !ok {
		return nil
	}
	return &v
}
