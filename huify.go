package xyebot

import (
	"math/rand"
	"regexp"
	"strings"
)

func huify(text string, gentle bool) string {
	huified := _huify(text)
	if huified == "" {
		return ""
	}
	if gentle {
		return Suggestions[rand.Intn(len(Suggestions))] + huified
	}
	return huified
}

func _huify(text string) string {
	const vowels string = "оеаяуюы"
	const rules_values string = "еяюи"
	rules := map[string]string{"о": "е", "а": "я", "у": "ю", "ы": "и"}
	NON_LETTERS, _ := regexp.Compile("[^а-яё-]+")
	ONLY_DASHES, _ := regexp.Compile("^-*$")
	PREFIX, _ := regexp.Compile("^[бвгджзйклмнпрстфхцчшщьъ]+")
	words := strings.Fields(text)
	if len(words) > 3 || len(words) < 1 {
		return ""
	}

	word := NON_LETTERS.ReplaceAllString(strings.ToLower(words[len(words)-1]), "")

	// Отдельная обработка слова бот
	if word == "бот" {
		return "хуебот"
	}
	// Пропускаем слова с дефисами, у которых после преобразования ничего, кроме дефисов не осталось
	if ONLY_DASHES.MatchString(word) {
		return ""
	}
	postfix := PREFIX.ReplaceAllString(word, "")
	// Пропускаем уже хуифицированные слова
	if len(postfix) < 6 || word[:4] == "ху" && strings.Index(rules_values, string(postfix[2:4])) >= 0 {
		return ""
	}
	// Пропускаем слова из стоп-листа
	if inStopList(word) {
		return ""
	}
	if _, ok := rules[postfix[0:2]]; ok {
		if strings.Index(vowels, postfix[2:4]) < 0 {
			return "ху" + rules[postfix[0:2]] + postfix[2:]
		} else {
			if huified, ok := rules[postfix[2:4]]; ok {
				return "ху" + huified + postfix[4:]
			} else {
				return "ху" + postfix[2:]
			}
		}
	} else {
		return "ху" + postfix
	}
}
