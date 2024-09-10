package xbytes

// IsWhitespace tests if a given byte is an ASCII whitespace character.
//
// Space characters are SPACE, TAB, LINE_FEED, and CARRIAGE_RETURN.
func IsWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// IsAlpha tests whether a given byte is an ASCII letter character.
func IsAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// IsNumeric tests whether a given byte is an ASCII number character.
func IsNumeric(b byte) bool {
	return b >= '0' && b <= '9'
}

// IsAlphaNumeric tests whether a given byte is either an ASCII letter or ASCII
// number character.
func IsAlphaNumeric(b byte) bool {
	return IsAlpha(b) || IsNumeric(b)
}

// IsWordChar tests whether a given byte is an ASCII letter, number, or
// underscore character.
//
// Aligns with the POSIX RegEx `[:word:]`
func IsWordChar(b byte) bool {
	return IsAlpha(b) || IsNumeric(b) || IsUnderscore(b)
}

// IsUnderscore tests whether the given byte is an ASCII underscore character.
func IsUnderscore(b byte) bool {
	return b == '_'
}
