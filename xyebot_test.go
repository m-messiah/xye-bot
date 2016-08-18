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
    AssertEqual(t, huify("привет"), "хуивет")
    AssertEqual(t, huify("привет бот"), "хуебот")
    AssertEqual(t, huify("доброе утро"), "хуютро")
    AssertEqual(t, huify("ты пьяный"), "хуяный")
    AssertEqual(t, huify("были"), "хуили")
    AssertEqual(t, huify("китайцы"), "хуитайцы")
}

func TestHuified(t *testing.T) {
    AssertEqual(t, huify("хуитайцы"), "")
    AssertEqual(t, huify("хуютро"), "")
    AssertEqual(t, huify("хутор"), "хуютор")
}
func TestNonRus(t *testing.T) {
    AssertEqual(t, huify("hello"), "")
    AssertEqual(t, huify("123"), "")
}
func TestDashed(t *testing.T) {
    AssertEqual(t, huify("когда-то"), "хуегда-то")
    AssertEqual(t, huify("semi-drive"), "")
    AssertEqual(t, huify("шалтай-болтай"), "хуялтай-болтай")
    AssertEqual(t, huify("https://www.edx.org/by-sec-li-mitx-3"), "")
}

func TestUrl(t *testing.T) {
    AssertEqual(t, huify("сайт.рф"), "хуяйтрф")
    AssertEqual(t, huify("http://www.ru"), "")
}
