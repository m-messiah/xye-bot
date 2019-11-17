package main

import (
	"fmt"
	"testing"
)

type TestWordParams struct {
	word        string
	huifiedWord string
	isHuified   bool
}

func CheckWord(t *testing.T, params TestWordParams) {
	t.Run(params.word, func(t *testing.T) {
		result, ok := tryHuifyWord(params.word)
		if result != params.huifiedWord {
			t.Error(result, "!=", params.huifiedWord)
		}
		if ok != params.isHuified {
			t.Error("isHuified?", ok, "!=", params.isHuified)
		}
	})
}

type TestPhraseParams struct {
	phrase        string
	amount        int
	huifiedPhrase string
}

func CheckPhrase(t *testing.T, params TestPhraseParams) {
	testName := fmt.Sprintf("%s_%d", params.phrase, params.amount)
	t.Run(testName, func(t *testing.T) {
		result := tryHuify(params.phrase, params.amount)
		if result != params.huifiedPhrase {
			t.Error(result, "!=", params.huifiedPhrase)
		}
	})
}

func TestRussian(t *testing.T) {
	tests := []TestWordParams{
		{"привет", "хуивет", true},
		{"были", "хуили", true},
		{"хутор", "хуютор", true},
		{"хахаха", "хахаха", false},
		{"ахаха", "ахаха", false},
		{"ахах", "ахах", false},
		{"хах", "хах", false},
		{"хха", "хха", false},
		{"аах", "аах", false},
		{"ах", "ах", false},
		{"в", "в", false},
		{"в ", "в", false},
		{" в", "в", false},
		{"", "", false},
	}

	for _, test := range tests {
		CheckWord(t, test)
	}
}

func TestUkrainian(t *testing.T) {
	tests := []TestWordParams{
		{"привіт", "хуівіт", true},
		{"вірила", "хуїрила", true},
	}

	for _, test := range tests {
		CheckWord(t, test)
	}
}

func TestHuified(t *testing.T) {
	tests := []TestWordParams{
		{"хуивет", "хуивет", false},
		{"хуютро", "хуютро", false},
	}

	for _, test := range tests {
		CheckWord(t, test)
	}
}

func TestSeveralWords(t *testing.T) {
	tests := []TestPhraseParams{
		{"привет", 5, "хуивет"},
		{"привет бот", 0, ""},
		{"привет бот", 1, "хуебот"},
		{"доброе утро", 1, "хуютро"},
		{"ты пьяный", 1, "хуяный"},
		{"привет бот", 1, "хуебот"},
		{"доброе утро", 4, "хуеброе хуютро"},
		{"доброе утро", 2, "хуеброе хуютро"},
		{"Мороз и солнце - день чудесный", 1, ""},
		{"Мороз и солнце - день чудесный", 2, "хуень хуюдесный"},
		{"Мороз и солнце - день чудесный", 3, "- хуень хуюдесный"},
		{"Мороз и солнце - день чудесный", 4, "хуелнце - хуень хуюдесный"},
		{"Мороз и солнце - день чудесный", 5, "и хуелнце - хуень хуюдесный"},
		{"Мороз и солнце - день чудесный", 6, "хуероз и хуелнце - хуень хуюдесный"},
		{"Выйду ночью в поле с конем", 10, "хуийду хуечью в хуеле с хуенем"},
		{"Выйду ночью в поле с конем", 1, ""},
		{"Было или не было, прошло или нет", 10, "хуило хуили не хуило хуешло хуили нет"},
		{"А не спеть ли мне песню?", 10, "а не хуеть ли мне хуесню"},
		{"А не и да", 10, ""},
	}

	for _, test := range tests {
		CheckPhrase(t, test)
	}
}

func TestNonRus(t *testing.T) {
	tests := []TestPhraseParams{
		{"hello", 2, ""},
		{"h", 2, ""},
		{"h w", 2, ""},
		{"h ", 2, ""},
		{" h", 2, ""},
		{"123", 2, ""},
	}

	for _, test := range tests {
		CheckPhrase(t, test)
	}
}
func TestDashed(t *testing.T) {
	tests := []TestPhraseParams{
		{"когда-то", 1, "хуегда-то"},
		{"шалтай-болтай", 1, "хуялтай-болтай"},
		{"semi-drive", 1, ""},
		{"https://www.edx.org/by-sec-li-mitx-3", 1, ""},
	}

	for _, test := range tests {
		CheckPhrase(t, test)
	}
}

func TestUrl(t *testing.T) {
	tests := []TestPhraseParams{
		{"сайт.рф", 1, "хуяйтрф"},
		{"http://www.ru", 1, ""},
	}

	for _, test := range tests {
		CheckPhrase(t, test)
	}
}
