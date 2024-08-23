package xbytes

func IsWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func IsAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func IsNumeric(b byte) bool {
	return b >= '0' && b <= '9'
}

func IsAlphaNumeric(b byte) bool {
	return IsAlpha(b) || IsNumeric(b)
}

func IsWordChar(b byte) bool {
	return IsAlpha(b) || IsNumeric(b) || IsUnderscore(b)
}

func IsUnderscore(b byte) bool {
	return b == '_'
}
