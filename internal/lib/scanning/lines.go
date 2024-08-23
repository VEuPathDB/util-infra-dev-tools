package scanning

import "vpdb-dev-tool/internal/lib/xbytes"

func FirstNonWhitespace(line []byte) (b byte, pos int) {
	for pos, b = range line {
		if !xbytes.IsWhitespace(b) {
			return
		}
	}

	return 0, -1
}
