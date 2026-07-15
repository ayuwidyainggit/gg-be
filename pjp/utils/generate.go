package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateCode(n int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var charsets = []rune("0123456789")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[rand.Intn(len(charsets))]
	}

	randomString := string(letters)
	randomInt, _ := strconv.Atoi(randomString)
	return randomInt
}
