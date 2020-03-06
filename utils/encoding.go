package utils

import "encoding/base64"

// Encode encrypts the content in base64 format
func Encode(msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(msg))
}

// Decode decrypts the base64 content to string
func Decode(msg string) string {
	decoded, _ := base64.StdEncoding.DecodeString(msg)
	return string(decoded)
}
