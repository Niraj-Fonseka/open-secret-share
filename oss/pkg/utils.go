package pkg

import (
	"math/rand"
	"strings"
	"time"
)

type Utils struct {
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewUtils() *Utils {
	return &Utils{}
}

func (u *Utils) GenerateUniqueKey() string {
	n := 20
	rand.Seed(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
