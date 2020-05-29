package utils

import (
	"math/rand"
	"time"
)

var letters = []byte("0123456789abcdefghijklmhopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numLetters =[]byte("0123456789")

// Generate a random string of n bits
func RandString(n uint) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func RandNum(n uint) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = numLetters[rand.Int63()%int64(len(letters))]
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
