package env

import "vpdb-dev-tool/internal/lib/xbytes"

// ResemblesEnvKey tests whether a given string could potentially be a valid
// environment variable name.
//
// Valid environment variable names must begin with a letter or underscore and
// must then consist entirely of letters, numbers, and/or underscores.
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
