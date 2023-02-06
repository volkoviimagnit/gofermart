package helpers

import (
	"math/rand"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
)

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomOrderNumber() string {
	return goluhn.Generate(16)
}

func RandomDigits(n int) string {
	rand.Seed(time.Now().UnixNano())

	letterRunes := []rune("1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
