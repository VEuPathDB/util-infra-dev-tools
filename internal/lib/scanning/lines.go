package scanning

import "vpdb-dev-tool/internal/lib/xbytes"

// FirstNonWhitespace returns the value and position of the first byte in the
// input that is not a whitespace character.
//
// If the input does not contain a non-whitespace character, this function
// returns a NULL byte and the position -1
func FirstNonWhitespace(line []byte) (b byte, pos int) {
	for pos, b = range line {
		if !xbytes.IsWhitespace(b) {
			return
		}
	}

	return 0, -1
}
