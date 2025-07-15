package utilities

import (
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	const CHARSET = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	for i := range bytes {
		bytes[i] = CHARSET[seededRand.Intn(len(CHARSET))]
	}

	return string(bytes)
}
