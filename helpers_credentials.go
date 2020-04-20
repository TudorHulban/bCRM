package main

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomString To be used for password salt.
func GenerateRandomString(targetLength int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]string, targetLength)

	randInt := func(min, max int) int {
		return min + rand.Intn(max-min)
	}
	for k := range result {
		result[k] = string(byte(randInt(65, 90)))
	}
	return strings.Join(result, "")
}

// HashPassword Hashes provided password with provided salt.
func HashPassword(pass, salt string) (string, error) {
	bytes, errHash := bcrypt.GenerateFromPassword([]byte(pass+salt), 14)
	return string(bytes), errHash
}
