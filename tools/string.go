package tools

import (
	"strconv"
	"time"
)

func GetCurrntTime() time.Time {
	return time.Now()
}

func StringToInt(e string) (int, error) {
	return strconv.Atoi(e)
}

func GetCurrntTimeStr() string {
	return time.Now().Format("2006/01/02 15:04:05")
}
