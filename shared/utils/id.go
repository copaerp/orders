package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateDisplayID() string {
	now := time.Now().UnixMilli()
	base36 := strings.ToUpper(strconv.FormatInt(now, 36))
	randDigit := strconv.FormatInt(int64(rand.Intn(36)), 36)
	return base36 + strings.ToUpper(randDigit)
}
