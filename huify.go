package xyebot

import (
	"math/rand"
	"regexp"
	"strings"
)

func huify(text string, gentle bool, amount int) string {
	huified := _huify(text, amount)
	if huified == "" {
		return ""
	}
	if gentle {
		return suggestions[rand.Intn(len(suggestions))] + huified
	}
	return huified
}

func _huify(text string, amount int) string {
	words := strings.Fields(text)
	// if len(words) > 3 || len(words) < 1 {
	if len(words) < 1 || len(words)-amount > 4 {
		return ""
	}
	var answer []string
	candidate_words := words
	if len(words) > amount {
		candidate_words = words[len(words)-amount:]
	}
	isHuified := false
	for _, word := range candidate_words {
		output, ok := _huify_word(word)
		if len(output) > 0 {
			answer = append(answer, output)
		}
		isHuified = isHuified || ok
	}
	if isHuified {
		return strings.Join(answer, " ")
	} else {
		return ""
	}
}

func _huify_word(text string) (string, bool) {
	const vowels string = "оеаяуюы"
	const rulesValues string = "еяюи"
	rules := map[string]string{"о": "е", "а": "я", "у": "ю", "ы": "и"}
	nonLetters, _ := regexp.Compile("[^а-яё-]+")
	onlyDashes, _ := regexp.Compile("^-*$")
	PREFIX, _ := regexp.Compile("^[бвгджзйклмнпрстфхцчшщьъ]+")

	word := nonLetters.ReplaceAllString(strings.ToLower(text), "")

	// Отдельная обработка слова бот
	if word == "бот" {
		return "хуебот", true
	}
	// Пропускаем слова с дефисами, у которых после преобразования ничего, кроме дефисов не осталось
	if onlyDashes.MatchString(word) {
		return word, false
	}
	postfix := PREFIX.ReplaceAllString(word, "")
	// Пропускаем уже хуифицированные слова
	if len(postfix) < 6 || word[:4] == "ху" && strings.Index(rulesValues, string(postfix[2:4])) >= 0 {
		return word, false
	}
	// Пропускаем слова из стоп-листа
	if inStopList(word) {
		return word, false
	}
	if _, ok := rules[postfix[0:2]]; ok {
		if strings.Index(vowels, postfix[2:4]) < 0 {
			return "ху" + rules[postfix[0:2]] + postfix[2:], true
		}
		if huified, ok := rules[postfix[2:4]]; ok {
			return "ху" + huified + postfix[4:], true
		}
		return "ху" + postfix[2:], true
	}
	return "ху" + postfix, true
}
