package xyebot

import (
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestRussian(c *C) {
	c.Check(_huify_word("привет"), Equals, "хуивет")
	c.Check(_huify_word("были"), Equals, "хуили")
	c.Check(_huify_word("китайцы"), Equals, "хуитайцы")
	c.Check(_huify_word("хахаха"), Equals, "")
	c.Check(_huify_word("ахаха"), Equals, "")
	c.Check(_huify_word("ахах"), Equals, "")
	c.Check(_huify_word("хах"), Equals, "")
	c.Check(_huify_word("хха"), Equals, "")
	c.Check(_huify_word("аах"), Equals, "")
	c.Check(_huify_word("ах"), Equals, "")
	c.Check(_huify_word("в"), Equals, "")
	c.Check(_huify_word("в "), Equals, "")
	c.Check(_huify_word(" в"), Equals, "")
	c.Check(_huify_word(""), Equals, "")
}

func (s *TestSuite) TestHuified(c *C) {
	c.Check(_huify_word("хуитайцы"), Equals, "")
	c.Check(_huify_word("хуютро"), Equals, "")
	c.Check(_huify_word("хутор"), Equals, "хуютор")
}

func (s *TestSuite) TestSeveralWords(c *C) {
	c.Check(_huify("привет", 5), Equals, "хуивет")
	c.Check(_huify("привет бот", 1), Equals, "хуебот")
	c.Check(_huify("доброе утро", 1), Equals, "хуютро")
	c.Check(_huify("ты пьяный", 1), Equals, "хуяный")
	c.Check(_huify("привет бот", 1), Equals, "хуебот")
	c.Check(_huify("доброе утро", 4), Equals, "хуеброе хуютро")
	c.Check(_huify("доброе утро", 2), Equals, "хуеброе хуютро")
	c.Check(_huify("Мороз и солнце - день чудесный", 1), Equals, "") // Слишком много слов
	c.Check(_huify("Мороз и солнце - день чудесный", 2), Equals, "хуень хуюдесный")
	c.Check(_huify("Мороз и солнце - день чудесный", 3), Equals, "хуень хуюдесный")
	c.Check(_huify("Мороз и солнце - день чудесный", 4), Equals, "хуелнце хуень хуюдесный")
	c.Check(_huify("Мороз и солнце - день чудесный", 5), Equals, "хуелнце хуень хуюдесный")
	c.Check(_huify("Мороз и солнце - день чудесный", 6), Equals, "хуероз хуелнце хуень хуюдесный")
}

func (s *TestSuite) TestNonRus(c *C) {
	c.Check(_huify("hello", 2), Equals, "")
	c.Check(_huify("h", 2), Equals, "")
	c.Check(_huify("h w", 2), Equals, "")
	c.Check(_huify("h ", 2), Equals, "")
	c.Check(_huify(" h", 2), Equals, "")
	c.Check(_huify("123", 2), Equals, "")
}
func (s *TestSuite) TestDashed(c *C) {
	c.Check(_huify("когда-то", 1), Equals, "хуегда-то")
	c.Check(_huify("semi-drive", 1), Equals, "")
	c.Check(_huify("шалтай-болтай", 1), Equals, "хуялтай-болтай")
	c.Check(_huify("https://www.edx.org/by-sec-li-mitx-3", 1), Equals, "")
}

func (s *TestSuite) TestUrl(c *C) {
	c.Check(_huify("сайт.рф", 1), Equals, "хуяйтрф")
	c.Check(_huify("http://www.ru", 1), Equals, "")
}
