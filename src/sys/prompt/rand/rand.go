//go:build !test

package rand

import (
	"math/rand"
	"strconv"
	"time"
)

func ThreeDigits() string {
	unix := time.Now().UnixNano()
	r := 100 + rand.New(rand.NewSource(unix)).Intn(99)
	return strconv.Itoa(r)
}
