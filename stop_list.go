package xyebot

import (
	"regexp"
)

func inStopList(word string) bool {
	if xaxaxa(word) {
		return true
	}
	return false
}

func xaxaxa(word string) bool {
	reHa, _ := regexp.Compile("^[ах]+$")
	return reHa.MatchString(word)
}
