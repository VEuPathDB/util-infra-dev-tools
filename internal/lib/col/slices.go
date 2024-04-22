package col

func SliceDeleteFirst[T comparable](value T, slice []T) ([]T, bool) {
	idx := -1
	for i := range slice {
		if slice[i] == value {
			idx = i
			break
		}
	}

	switch {
	case idx < 0:
		return slice, false
	case idx == 0:
		return slice[1:], true
	case idx == len(slice)-1:
		return slice[:len(slice)-1], true
	}

	out := make([]T, len(slice)-1)
	copy(out, slice[:idx])
	copy(out[idx:], slice[idx+1:])

	return out, true
}
