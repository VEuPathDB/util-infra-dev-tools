package env

import "vpdb-dev-tool/internal/lib/xbytes"

func ResemblesEnvKey(key string) bool {
	if len(key) < 1 {
		return false
	}

	if !xbytes.IsAlpha(key[0]) && !xbytes.IsUnderscore(key[0]) {
		return false
	}

	for i := 1; i < len(key); i++ {
		if !xbytes.IsWordChar(key[i]) {
			return false
		}
	}

	return true
}
