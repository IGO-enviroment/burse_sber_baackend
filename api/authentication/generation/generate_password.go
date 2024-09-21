package generation

import (
	"crypto/rand"
	"math/big"
)

const passwordLength = 16
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GeneratePassword() string {
	var result string

	for i := 0; i < passwordLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		result += string(charset[randomIndex.Int64()])
	}

	return result
}
