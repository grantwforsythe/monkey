package strings

// Determine if a character is whitespace or not
func IsWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// Determine if a charater is a letter or not
func IsLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Determine if a charater is a digit or not
func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
