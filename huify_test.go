package xyebot

import (
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestRussian(c *C) {
	result, ok := TryHuifyWord("привет")
	c.Check(result, Equals, "хуивет")
	c.Check(ok, Equals, true)
	result, ok = TryHuifyWord("были")
	c.Check(result, Equals, "хуили")
	c.Check(ok, Equals, true)
	result, ok = TryHuifyWord("китайцы")
	c.Check(result, Equals, "хуитайцы")
	c.Check(ok, Equals, true)
	result, ok = TryHuifyWord("хахаха")
	c.Check(result, Equals, "хахаха")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("ахаха")
	c.Check(result, Equals, "ахаха")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("ахах")
	c.Check(result, Equals, "ахах")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("хах")
	c.Check(result, Equals, "хах")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("хха")
	c.Check(result, Equals, "хха")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("аах")
	c.Check(result, Equals, "аах")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("ах")
	c.Check(result, Equals, "ах")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("в")
	c.Check(result, Equals, "в")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("в ")
	c.Check(result, Equals, "в")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord(" в")
	c.Check(result, Equals, "в")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("")
	c.Check(result, Equals, "")
	c.Check(ok, Equals, false)
}

func (s *TestSuite) TestHuified(c *C) {
	result, ok := TryHuifyWord("хуитайцы")
	c.Check(result, Equals, "хуитайцы")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("хуютро")
	c.Check(result, Equals, "хуютро")
	c.Check(ok, Equals, false)
	result, ok = TryHuifyWord("хутор")
	c.Check(result, Equals, "хуютор")
	c.Check(ok, Equals, true)
}

func (s *TestSuite) TestSeveralWords(c *C) {
	c.Check(TryHuify("привет", 5), Equals, "хуивет")
	c.Check(TryHuify("привет бот", 0), Equals, "")
	c.Check(TryHuify("привет бот", 1), Equals, "хуебот")
	c.Check(TryHuify("доброе утро", 1), Equals, "хуютро")
	c.Check(TryHuify("ты пьяный", 1), Equals, "хуяный")
	c.Check(TryHuify("привет бот", 1), Equals, "хуебот")
	c.Check(TryHuify("доброе утро", 4), Equals, "хуеброе хуютро")
	c.Check(TryHuify("доброе утро", 2), Equals, "хуеброе хуютро")
	c.Check(TryHuify("Мороз и солнце - день чудесный", 1), Equals, "") // Слишком много слов
	c.Check(TryHuify("Мороз и солнце - день чудесный", 2), Equals, "хуень хуюдесный")
	c.Check(TryHuify("Мороз и солнце - день чудесный", 3), Equals, "- хуень хуюдесный")
	c.Check(TryHuify("Мороз и солнце - день чудесный", 4), Equals, "хуелнце - хуень хуюдесный")
	c.Check(TryHuify("Мороз и солнце - день чудесный", 5), Equals, "и хуелнце - хуень хуюдесный")
	c.Check(TryHuify("Мороз и солнце - день чудесный", 6), Equals, "хуероз и хуелнце - хуень хуюдесный")
	c.Check(TryHuify("Выйду ночью в поле с конем", 10), Equals, "хуийду хуечью в хуеле с хуенем")
	c.Check(TryHuify("Выйду ночью в поле с конем", 1), Equals, "")
	c.Check(TryHuify("Было или не было, прошло или нет", 10), Equals, "хуило хуили не хуило хуешло хуили нет")
	c.Check(TryHuify("А не спеть ли мне песню?", 10), Equals, "а не хуеть ли мне хуесню")
	c.Check(TryHuify("А не и да", 10), Equals, "")
}

func (s *TestSuite) TestNonRus(c *C) {
	c.Check(TryHuify("hello", 2), Equals, "")
	c.Check(TryHuify("h", 2), Equals, "")
	c.Check(TryHuify("h w", 2), Equals, "")
	c.Check(TryHuify("h ", 2), Equals, "")
	c.Check(TryHuify(" h", 2), Equals, "")
	c.Check(TryHuify("123", 2), Equals, "")
}
func (s *TestSuite) TestDashed(c *C) {
	c.Check(TryHuify("когда-то", 1), Equals, "хуегда-то")
	c.Check(TryHuify("semi-drive", 1), Equals, "")
	c.Check(TryHuify("шалтай-болтай", 1), Equals, "хуялтай-болтай")
	c.Check(TryHuify("https://www.edx.org/by-sec-li-mitx-3", 1), Equals, "")
}

func (s *TestSuite) TestUrl(c *C) {
	c.Check(TryHuify("сайт.рф", 1), Equals, "хуяйтрф")
	c.Check(TryHuify("http://www.ru", 1), Equals, "")
}
