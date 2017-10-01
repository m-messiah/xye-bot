package xyebot

import (
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestHumanNumeral(c *C) {
	suffixes := [3]string{"слово", "слова", "слов"}
	format_string := "Здесь %s"
	c.Check(humanNumeral(format_string, 1, suffixes), Equals, "Здесь 1 слово")
	c.Check(humanNumeral(format_string, 2, suffixes), Equals, "Здесь 2 слова")
	c.Check(humanNumeral(format_string, 3, suffixes), Equals, "Здесь 3 слова")
	c.Check(humanNumeral(format_string, 4, suffixes), Equals, "Здесь 4 слова")
	c.Check(humanNumeral(format_string, 5, suffixes), Equals, "Здесь 5 слов")
	c.Check(humanNumeral(format_string, 8, suffixes), Equals, "Здесь 8 слов")
	c.Check(humanNumeral(format_string, 10, suffixes), Equals, "Здесь 10 слов")
	c.Check(humanNumeral(format_string, 11, suffixes), Equals, "Здесь 11 слов")
	c.Check(humanNumeral(format_string, 14, suffixes), Equals, "Здесь 14 слов")
	c.Check(humanNumeral(format_string, 18, suffixes), Equals, "Здесь 18 слов")
	c.Check(humanNumeral(format_string, 21, suffixes), Equals, "Здесь 21 слово")
	c.Check(humanNumeral(format_string, 22, suffixes), Equals, "Здесь 22 слова")
	c.Check(humanNumeral(format_string, 39, suffixes), Equals, "Здесь 39 слов")
}
