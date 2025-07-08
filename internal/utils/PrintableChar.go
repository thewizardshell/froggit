package utils

// IsPrintableChar checks whether a rune is a printable ASCII or extended character.
func IsPrintableChar(r rune) bool {
	return (r >= 32 && r <= 126) || (r >= 128 && r <= 255) // Extended ASCII range
}
