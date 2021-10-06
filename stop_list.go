package main

import (
	"regexp"
)

var xaxaxaRe = regexp.MustCompile("^[ах]+$")

func inStopList(word string) bool {
	return isXAXAXA(word)
}

func isXAXAXA(word string) bool {
	return xaxaxaRe.MatchString(word)
}
