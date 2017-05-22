package xyebot

import (
	"testing"
)

func AssertEqual(t *testing.T, got, expected string) {
	if got != expected {
		t.Error("Expected " + expected + ", got " + got)
	}
}

func TestRussian(t *testing.T) {
	AssertEqual(t, _huify("привет"), "хуивет")
	AssertEqual(t, _huify("привет бот"), "хуебот")
	AssertEqual(t, _huify("доброе утро"), "хуютро")
	AssertEqual(t, _huify("ты пьяный"), "хуяный")
	AssertEqual(t, _huify("были"), "хуили")
	AssertEqual(t, _huify("китайцы"), "хуитайцы")
	AssertEqual(t, _huify("хахаха"), "")
	AssertEqual(t, _huify("ахаха"), "")
	AssertEqual(t, _huify("ахах"), "")
	AssertEqual(t, _huify("хах"), "")
	AssertEqual(t, _huify("хха"), "")
	AssertEqual(t, _huify("аах"), "")
	AssertEqual(t, _huify("ах"), "")
	AssertEqual(t, _huify("в"), "")
	AssertEqual(t, _huify("в "), "")
	AssertEqual(t, _huify(" в"), "")
	AssertEqual(t, _huify(""), "")
}

func TestHuified(t *testing.T) {
	AssertEqual(t, _huify("хуитайцы"), "")
	AssertEqual(t, _huify("хуютро"), "")
	AssertEqual(t, _huify("хутор"), "хуютор")
}
func TestNonRus(t *testing.T) {
	AssertEqual(t, _huify("hello"), "")
	AssertEqual(t, _huify("h"), "")
	AssertEqual(t, _huify("h w"), "")
	AssertEqual(t, _huify("h "), "")
	AssertEqual(t, _huify(" h"), "")
	AssertEqual(t, _huify("123"), "")
}
func TestDashed(t *testing.T) {
	AssertEqual(t, _huify("когда-то"), "хуегда-то")
	AssertEqual(t, _huify("semi-drive"), "")
	AssertEqual(t, _huify("шалтай-болтай"), "хуялтай-болтай")
	AssertEqual(t, _huify("https://www.edx.org/by-sec-li-mitx-3"), "")
}

func TestUrl(t *testing.T) {
	AssertEqual(t, _huify("сайт.рф"), "хуяйтрф")
	AssertEqual(t, _huify("http://www.ru"), "")
}
