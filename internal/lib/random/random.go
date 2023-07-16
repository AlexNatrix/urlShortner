package random

import (
	"math/rand"
	"time"
)

func NewRandomStr(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, size)
	for i := range b {
		b[i] = rune(65 + rnd.Intn(200)%26)
	}
	return string(b)
}
