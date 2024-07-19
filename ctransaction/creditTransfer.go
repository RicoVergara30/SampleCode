package ctransaction

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Iftgenerate() string {
	uniqueDigit := rand.Intn(1000)
	uniqueDigit1 := rand.Int63n(10000000000000) // Use Int63n for larger range
	return fmt.Sprintf("IFT%v-%v", uniqueDigit, uniqueDigit1)
}

// func GenerateInstructionID(max_digits, current_count int) string {
// 	max_digits = max_digits - 1
// 	var instructionID string
// 	current_length := len(strconv.Itoa(current_count))

// 	if current_length <= max_digits {
// 		current_count++
// 		for strL := 0; strL <= max_digits-current_length; strL++ {
// 			instructionID += "0"
// 		}
// 	} else {
// 		current_count = 1
// 		for strL := 0; strL <= max_digits-current_length; strL++ {
// 			instructionID += "0"
// 		}
// 	}

// 	instructionID += strconv.Itoa(current_count)
// 	return instructionID
// }

// func GenerateInstructionID(uniqueDigit int) string {
// 	// Initialize random seed based on current time
// 	rand.Seed(time.Now().UnixNano())

// 	// Generate a random number
// 	randomNumber := rand.Intn(1000000) // 6 digits random number

//		// Format the ID
//		return fmt.Sprintf("20240219CAMZPHM2%06dB0000000000%d", randomNumber, uniqueDigit)
//	}
// func GenerateInstructionID(len int) string {
// 	now := time.Now()
// 	nowTime := now.Unix()
// 	var identifier string

// 	for ctr := 0; ctr < len; ctr++ {
// 		rndNum := strconv.Itoa(rand.Intn(9-0) + 1)
// 		identifier += rndNum
// 	}

// 	identifier = strconv.Itoa(int(nowTime)) + identifier
// 	return identifier
// }

func GenerateInstructionID() string {
	// Current date in YYYYMMDD format
	datePart := time.Now().Format("20060102")

	// Static bank identifier code
	bankIdentifierCode := "CAMZPHM2XXXB"

	// Generate a 14-digit random number
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < 14; i++ {
		sb.WriteByte(byte('0' + rand.Intn(10)))
	}

	uniqueNumber := sb.String()

	return fmt.Sprintf("%s%s%s", datePart, bankIdentifierCode, uniqueNumber)
}
