package binary

func SearchL[T any](n int, f func(int) (bool, T, error)) (int, T, error) {
	var zero, elLeft T
	l, r := -1, n
	for r-l > 1 {
		m := int(uint(r+l) >> 1)
		ok, current, err := f(m)
		if err != nil {
			return -1, zero, err
		}
		if ok {
			l = m
			elLeft = current
		} else {
			r = m
		}
	}
	return l, elLeft, nil
}

func SearchR[T any](n int, f func(int) (bool, T, error)) (int, T, error) {
	var zero, elRight T
	l, r := -1, n
	for r-l > 1 {
		m := int(uint(r+l) >> 1)
		ok, current, err := f(m)
		if err != nil {
			return -1, zero, err
		}
		if ok {
			l = m
		} else {
			r = m
			elRight = current
		}
	}
	return r, elRight, nil
}

func SearchWithError(n int, f func(int) (bool, error)) (int, error) {
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1)
		ok, err := f(h)
		if err != nil {
			return -1, err
		}
		if !ok {
			i = h + 1
		} else {
			j = h
		}
	}
	return i, nil
}
