package utilities

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {

	charset := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return fmt.Sprintf("%s", string(b))
}
