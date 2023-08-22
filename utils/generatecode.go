package utils

import (
	"crypto/rand"
)

const (
	uniqueCodeLength = 6
	chars            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateUniqueCode generates a unique alphanumeric code with 6 digits
func GenerateUniqueCode() string {
	code := make([]byte, uniqueCodeLength)
	rand.Read(code)

	for i := 0; i < uniqueCodeLength; i++ {
		code[i] = chars[code[i]%byte(len(chars))]
	}

	return string(code)
}
