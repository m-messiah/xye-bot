package xyebot

import (
	"encoding/json"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

func huify(text string) string {
	const vowels string = "оеаяуюы"
	const rules_values string = "еяюи"
	rules := map[string]string{"о": "е", "а": "я", "у": "ю", "ы": "и"}
	NON_LETTERS, _ := regexp.Compile("[^а-яё-]+")
	ONLY_DASHES, _ := regexp.Compile("^-*$")
	PREFIX, _ := regexp.Compile("^[бвгджзйклмнпрстфхцчшщьъ]+")
	words := strings.Fields(text)
	if len(words) > 3 {
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
	if word[:4] == "ху" && strings.Index(rules_values, string(postfix[2:4])) >= 0 || len(postfix) < 6 {
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

func init() {
	DELAY := make(map[int64]int)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)

		var update tgbotapi.Update
		json.Unmarshal(bytes, &update)
		if update.Message == nil {
			return
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if strings.Contains(update.Message.Text, "/start") || strings.Contains(update.Message.Text, "/help") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот-хуебот.\nЯ буду хуифицировать некоторые из твоих фраз")
			json.NewEncoder(w).Encode(msg)
			return
		} else {
			if _, ok := DELAY[update.Message.Chat.ID]; ok {
				DELAY[update.Message.Chat.ID] -= 1
			} else {
				DELAY[update.Message.Chat.ID] = rand.Intn(4)
			}
			if DELAY[update.Message.Chat.ID] == 0 {
				delete(DELAY, update.Message.Chat.ID)
				output := huify(update.Message.Text)
				if output != "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, output)
					json.NewEncoder(w).Encode(msg)
					return
				}
			}
		}
	})
}
