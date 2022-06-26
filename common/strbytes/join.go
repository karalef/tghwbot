package strbytes

// Join joins strings.
func Join(s []string, d byte) string {
	switch len(s) {
	case 0:
		return ""
	case 1:
		return s[0]
	}
	n := len(s) - 1
	for i := 0; i < len(s); i++ {
		n += len(s[i])
	}

	b := make([]byte, 0, n)
	b = append(b, s[0]...)
	for i := 1; i < len(s); i++ {
		b = append(b, d)
		b = append(b, s[i]...)
	}

	return B2s(b)
}
