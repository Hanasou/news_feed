package util

func FilterChar(str string, filter rune) string {
	result := ""
	for _, c := range str {
		if c != filter {
			result += string(c)
		}
	}
	return result
}
