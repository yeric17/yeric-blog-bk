package utils

import "math/rand"

func RandomString(length int) string {
	stringRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)
	for i := range b {
		b[i] = stringRunes[rand.Intn(len(stringRunes))]
	}
	return string(b)
}
