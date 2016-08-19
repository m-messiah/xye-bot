package xyebot

import (
    "math/rand"
    "regexp"
    "strings"
)

var Suggestions = [...]string{
    "Возможно, вы имели ввиду: ",
    "Извините, ",
    "Извините, но ",
    "Простите, ",
    "Подождите, ",
}

func huify(text string) string {
    // Cultural huify
    huified := _huify(text)
    if huified == "" {
        return ""
    }
    if rand.Intn(100) < 30 {
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

    if word == "бот" {
        return "хуебот"
    }
    if ONLY_DASHES.MatchString(word) {
        return ""
    }
    postfix := PREFIX.ReplaceAllString(word, "")
    if len(postfix) < 6 || word[:4] == "ху" && strings.Index(rules_values, string(postfix[2:4])) >= 0 {
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
