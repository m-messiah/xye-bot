package xyebot

import (
	"regexp"
)

var xaxaxaRe = regexp.MustCompile("^[ах]+$")

func inStopList(word string) bool {
	if isXAXAXA(word) {
		return true
	}
	return false
}

func isXAXAXA(word string) bool {
	return xaxaxaRe.MatchString(word)
}
