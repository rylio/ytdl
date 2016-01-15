package ytdl

import (
	"math/rand"
	"strings"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func randString(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

func TestReverseString(t *testing.T) {
	strs, reverseStrs := make([]string, 20), make([]string, 20)

	for i := 0; i < len(strs); i++ {
		strs[i] = randString(rand.Intn(50))
		s := strings.Split(strs[i], "")
		reverseStringSlice(s)
		reverseStrs[i] = strings.Join(s, "")
	}

	for i, s := range strs {
		rs := reverseStrs[i]
		for j, c := range s {
			if c != rune(rs[len(rs)-1-j]) {
				t.Fail()
			}
		}
	}
}

func TestInterfaceToString(t *testing.T) {

	values := map[interface{}]string{
		8:         "8",
		int64(16): "16",
		"hello":   "hello",
		0.01:      "0.01",
	}

	for k, v := range values {
		if interfaceToString(k) != v {
			t.Error("Value:", v, " != ", v)
		}
	}

}
