package types

// GoResult is a wrapper for a goroutine result with error handling
type GoResult[T any] struct {
	Response T
	Err      error
}

// Done is an alias for bool that indicates whether a goroutine is done.
// This is to have clarity between a bool result and a done indicator.
type Done = bool

// Tuple is a wrapper for a tuple with 2 items
type Tuple[T any, U any] struct {
	Item1 T
	Item2 U
}

// Tuple3 is a wrapper for a tuple with 3 items
type Tuple3[T any, U any, V any] struct {
	Item1 T
	Item2 U
	Item3 V
}

// Tuple4 is a wrapper for a tuple with 4 items
type Tuple4[T any, U any, V any, W any] struct {
	Item1 T
	Item2 U
	Item3 V
	Item4 W
}
