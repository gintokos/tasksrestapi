package id

import (
	"math/rand"
	"strconv"
	"time"
)

func ValidateID(id string) (int64, bool) {
	idint64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, false
	}

	return idint64, true
}

func GenerateRandomID() int64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Int63()
}
