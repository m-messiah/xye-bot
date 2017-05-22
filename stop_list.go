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
	re_ha, _ := regexp.Compile("^[ах]+$")
	return re_ha.MatchString(word)
}
