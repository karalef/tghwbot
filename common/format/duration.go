package format

// seconds
const (
	Nanosecond uint64 = 1

	Microsecond = 1000 * Nanosecond
	Millisecond = 1000 * Microsecond
	Second      = 1000 * Millisecond

	Minute = 60 * Second
	Hour   = 60 * Minute
)
