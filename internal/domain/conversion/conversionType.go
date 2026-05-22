package conversion

import (
	"fmt"
	"strconv"
)

func StringToUint(word string) uint {
	value, err := strconv.ParseUint(word, 10, 64)
	if err != nil {
		fmt.Println("")

	}
	return uint(value)
}

func Float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'f', 0, 64)
}

func UintToString(valor uint) string {
	return strconv.Itoa(int(valor))
}

func IntToString(valor int) string {
	return strconv.Itoa(int(valor))
}
