package main

import (
	"math/rand"
	"regexp"
	"strings"
)

var (
	prefixToSkipRe = regexp.MustCompile("^[бвгджзйклмнпрстфхцчшщьъ]+")
	nonLettersRe   = regexp.MustCompile("[^а-яёії-]+")
	onlyDashesRe   = regexp.MustCompile("^-*$")
	vowelsRules    = map[string]string{"о": "е", "а": "я", "у": "ю", "ы": "и"}
	vowelsRulesUA  = map[string]string{"о": "е", "а": "я", "у": "ю", "ы": "и", "и": "і", "і": "ї"}
)

const (
	applyUARules  string = "ії"
	rulesValues   string = "еяюиії"
	huifiedPrefix string = "ху"
	vowels        string = "оеаяуюыі"
)

// Huify given text by gentleness and limited amount
func Huify(text string, amount int) string {
	huified := tryHuify(text, amount)
	if huified == "" {
		return ""
	}
	return suggestions[rand.Intn(len(suggestions))] + huified
}

func tryHuify(text string, amount int) string {
	words := strings.Fields(text)
	if len(words) < 1 || len(words)-amount > 4 {
		return ""
	}
	var answer []string
	candidateWords := words
	if len(words) > amount {
		candidateWords = words[len(words)-amount:]
	}
	isHuified := false
	for _, word := range candidateWords {
		output, ok := tryHuifyWord(word)
		if len(output) > 0 {
			answer = append(answer, output)
		}
		isHuified = isHuified || ok
	}
	if isHuified {
		return strings.Join(answer, " ")
	}
	return ""
}

func isHuifyApplicable(word string) (*string, bool) {
	// Пропускаем слова с дефисами, у которых после преобразования ничего, кроме дефисов не осталось
	if onlyDashesRe.MatchString(word) {
		return nil, false
	}
	postfix := prefixToSkipRe.ReplaceAllString(word, "")
	// Пропускаем уже хуифицированные слова
	if len(postfix) < 6 || word[:4] == huifiedPrefix && strings.Contains(rulesValues, postfix[2:4]) {
		return nil, false
	}
	// Пропускаем слова из стоп-листа
	if inStopList(word) {
		return nil, false
	}

	return &postfix, true
}

func isUAWord(word string) bool {
	for _, letter := range applyUARules {
		if strings.Contains(word, string(letter)) {
			return true
		}
	}
	return false
}

func huifyWord(postfix string, rules map[string]string) string {
	if _, ok := rules[postfix[0:2]]; ok {
		if !strings.Contains(vowels, postfix[2:4]) {
			return huifiedPrefix + rules[postfix[0:2]] + postfix[2:]
		}
		if huified, ok := rules[postfix[2:4]]; ok {
			return huifiedPrefix + huified + postfix[4:]
		}
		return huifiedPrefix + postfix[2:]
	}
	return huifiedPrefix + postfix
}

func getRules(word string) map[string]string {
	if isUAWord(word) {
		return vowelsRulesUA
	}
	return vowelsRules
}

func tryHuifyWord(text string) (string, bool) {
	word := nonLettersRe.ReplaceAllString(strings.ToLower(text), "")

	// Отдельная обработка слова бот
	if word == "бот" {
		return "хуебот", true
	}

	if word == "путин" {
		return "хуйло", true
	}

	if postfix, ok := isHuifyApplicable(word); ok {
		return huifyWord(*postfix, getRules(word)), true
	}

	return word, false
}
