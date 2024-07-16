package ctransaction

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ReferenceId() string {
	uniqueDigit := rand.Intn(10000)
	return fmt.Sprintf("20240219CAMZPHM2XXXB0000000000%d", uniqueDigit)
}

func Iftgenerate() string {
	uniqueDigit := rand.Intn(1000)
	uniqueDigit1 := rand.Int63n(10000000000000) // Use Int63n for larger range
	return fmt.Sprintf("IFT%v-%v", uniqueDigit, uniqueDigit1)
}
