package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Hash takes a string and returns its md5 hashed value
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// VerifyMD5 checks if the raw password matches the hashed md5 password
func VerifyMD5(raw, hashed string) bool {
	return MD5Hash(raw) == hashed
}
