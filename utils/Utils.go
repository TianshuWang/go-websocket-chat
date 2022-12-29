package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int()
}
