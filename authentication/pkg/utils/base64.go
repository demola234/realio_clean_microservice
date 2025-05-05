package utils

import (
	"encoding/base64"
	"strings"
)

// Base64URLDecode decodes a base64url encoded string
func Base64URLDecode(str string) ([]byte, error) {
	// Replace URL-safe characters
	str = strings.Replace(str, "-", "+", -1)
	str = strings.Replace(str, "_", "/", -1)

	// Add padding if needed
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	return base64.StdEncoding.DecodeString(str)
}
