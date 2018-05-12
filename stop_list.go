package xyebot

import (
	"regexp"
)

var XAXAXA_RE = regexp.MustCompile("^[ах]+$")

func InStopList(word string) bool {
	if IsXAXAXA(word) {
		return true
	}
	return false
}

func IsXAXAXA(word string) bool {
	return XAXAXA_RE.MatchString(word)
}
