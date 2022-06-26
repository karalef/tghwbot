package common

// Flag is a small set of booleans.
type Flag uint8

// Set sets bits.
func (f *Flag) Set(flag Flag) {
	*f |= flag
}

// Clear clears bits.
func (f *Flag) Clear(flag Flag) {
	*f &^= flag
}

// Toggle toggles bits.
func (f *Flag) Toggle(flag Flag) {
	*f ^= flag
}

// Has true if flag contains flag bits.
func (f Flag) Has(flag Flag) bool {
	return f&flag != 0
}
