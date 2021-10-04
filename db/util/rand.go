package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(min, max int) int {
	return min + int(rand.Int63n(int64(max-min+1)))
}

func RandStr(n int) string {
	s := strings.Builder{}
	for i := 0; i < n; i++ {
		s.WriteByte('a' + byte(rand.Int63n(26)))
	}
	return s.String()
}

func RandDigits(n int) string {
	s := strings.Builder{}
	for i := 0; i < n; i++ {
		s.WriteByte('0' + byte(rand.Int63n(10)))
	}
	return s.String()
}

func RandType() string {
	types := []string{"salaried", "hourly", "commissioned"}
	return types[rand.Int63n(3)]
}
