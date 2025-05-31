package unscramble

func IReal(s string) string {
	var result string
	for len(s) > 50 {
		chunk := s[:50]
		s = s[50:]
		if len(s) < 2 {
			result += chunk
		} else {
			result += obfusc50(chunk)
		}
	}
	result += s
	return result
}

func obfusc50(s string) string {
	r := []rune(s)
	for i := 0; i < 5; i++ {
		r[i], r[49-i] = r[49-i], r[i]
	}
	for i := 10; i < 24; i++ {
		r[i], r[49-i] = r[49-i], r[i]
	}
	return string(r)
}
