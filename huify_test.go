package xyebot

import (
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestRussian(c *C) {
	result, ok := _huify_word("привет")
	c.Check(result, Equals, "хуивет")
	c.Check(ok, Equals, true)
	result, ok = _huify_word("были")
	c.Check(result, Equals, "хуили")
	c.Check(ok, Equals, true)
	result, ok = _huify_word("китайцы")
	c.Check(result, Equals, "хуитайцы")
	c.Check(ok, Equals, true)
	result, ok = _huify_word("хахаха")
	c.Check(result, Equals, "хахаха")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("ахаха")
	c.Check(result, Equals, "ахаха")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("ахах")
	c.Check(result, Equals, "ахах")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("хах")
	c.Check(result, Equals, "хах")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("хха")
	c.Check(result, Equals, "хха")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("аах")
	c.Check(result, Equals, "аах")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("ах")
	c.Check(result, Equals, "ах")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("в")
	c.Check(result, Equals, "в")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("в ")
	c.Check(result, Equals, "в")
	c.Check(ok, Equals, false)
	result, ok = _huify_word(" в")
	c.Check(result, Equals, "в")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("")
	c.Check(result, Equals, "")
	c.Check(ok, Equals, false)
}

func (s *TestSuite) TestHuified(c *C) {
	result, ok := _huify_word("хуитайцы")
	c.Check(result, Equals, "хуитайцы")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("хуютро")
	c.Check(result, Equals, "хуютро")
	c.Check(ok, Equals, false)
	result, ok = _huify_word("хутор")
	c.Check(result, Equals, "хуютор")
	c.Check(ok, Equals, true)
}

func (s *TestSuite) TestSeveralWords(c *C) {
	c.Check(_huify("привет", 5), Equals, "хуивет")
	c.Check(_huify("привет бот", 0), Equals, "")
	c.Check(_huify("привет бот", 1), Equals, "хуебот")
	c.Check(_huify("доброе утро", 1), Equals, "хуютро")
	c.Check(_huify("ты пьяный", 1), Equals, "хуяный")
	c.Check(_huify("привет бот", 1), Equals, "хуебот")
	c.Check(_huify("доброе утро", 4), Equals, "хуеброе хуютро")
	c.Check(_huify("доброе утро", 2), Equals, "хуеброе хуютро")
	c.Check(_huify("Мороз и солнце - день чудесный", 1), Equals, "") // Слишком много слов
	c.Check(_huify("Мороз и солнце - день чудесный", 2), Equals, "хуень хуюдесный")
	c.Check(_huify("Мороз и солнце - день чудесный", 3), Equals, "- хуень хуюдесный")
	c.Check(_huify("Мороз и солнце - день чудесный", 4), Equals, "хуелнце - хуень хуюдесный")
	c.Check(_huify("Мороз и солнце - день чудесный", 5), Equals, "и хуелнце - хуень хуюдесный")
	c.Check(_huify("Мороз и солнце - день чудесный", 6), Equals, "хуероз и хуелнце - хуень хуюдесный")
	c.Check(_huify("Выйду ночью в поле с конем", 10), Equals, "хуийду хуечью в хуеле с хуенем")
	c.Check(_huify("Выйду ночью в поле с конем", 1), Equals, "")
	c.Check(_huify("Было или не было, прошло или нет", 10), Equals, "хуило хуили не хуило хуешло хуили нет")
	c.Check(_huify("А не спеть ли мне песню?", 10), Equals, "а не хуеть ли мне хуесню")
	c.Check(_huify("А не и да", 10), Equals, "")
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
