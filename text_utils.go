package xyebot

import (
	"fmt"
	"strconv"
)

func humanNumeral(format_message string, num int, suffixes [3]string) string {
	numeral := strconv.Itoa(num) + " "
	if num > 5 && num < 21 {
		numeral += suffixes[2]
	} else {
		switch num % 10 {
		case 1:
			numeral += suffixes[0]
		case 2, 3, 4:
			numeral += suffixes[1]
		case 5, 6, 7, 8, 9, 0:
			numeral += suffixes[2]
		}
	}
	return fmt.Sprintf(format_message, numeral)
}
