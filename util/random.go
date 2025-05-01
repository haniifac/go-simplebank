package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var currencies = []string{"USD", "EUR", "CAD"}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		randomIdx := rand.Intn(len(alphabet))
		sb.WriteByte(alphabet[randomIdx])
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(10)
}

func RandomMoney() int64 {
	return int64(RandomInt(500, 1000))
}

func RandomCurrency() string {
	return currencies[rand.Intn(len(currencies))]
}
