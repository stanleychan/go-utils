package utils

import (
	"math/rand"
	"time"
)

var letters = []rune("0123456789abcdefghijklmhopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Generate a random string of n bits
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(letters[rand.Int63()%int64(len(letters))])
	}
	return string(b)
}

/** Generate a random bool value.
 *  Usage: rb := GetRandBool(); rb.Bool();
 */
type RandBool struct {
	src       rand.Source
	cache     int64
	remaining int
}

func (b *RandBool) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}
	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--
	return result
}

func GetRandBool() *RandBool {
	return &RandBool{
		src: rand.NewSource(time.Now().UnixNano()),
	}
}
