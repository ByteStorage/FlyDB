package randkv

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	randStr = rand.New(rand.NewSource(time.Now().Unix())) // 伪随机数
	letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

// GetTestKey Gets the key used by the test
func GetTestKey(i int) []byte {
	return []byte(fmt.Sprintf("flydb-key-%09d", i))
}

// RandomValue Generate random values for testing
func RandomValue(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[randStr.Intn(len(letters))]
	}
	return []byte("flydb-value-" + string(b))
}
