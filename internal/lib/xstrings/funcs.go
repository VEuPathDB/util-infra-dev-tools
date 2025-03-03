package xstrings

func Coalesce(first string, more ...string) string {
	if len(first) > 0 {
		return first
	}

	for _, s := range more {
		if len(s) > 0 {
			return s
		}
	}

	return ""
}
