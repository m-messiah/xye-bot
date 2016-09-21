package xyebot

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

func sendMessage(w http.ResponseWriter, chat_id int64, text string) {
	msg := Response{Chatid: chat_id, Text: text, Method: "sendMessage"}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

func init() {
	DELAY := make(map[int64]int)
	GENTLE := make(map[int64]bool)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		ctx := appengine.NewContext(r)

		var update Update
		json.Unmarshal(bytes, &update)
		if update.Message == nil {
			return
		}
		log.Debugf(ctx, string(bytes))
		if _, ok := GENTLE[update.Message.Chat.ID]; !ok {
			GENTLE[update.Message.Chat.ID] = true
		}

		if strings.Contains(update.Message.Text, "/start") {
			message := "Привет! Я бот-хуебот.\nЯ буду хуифицировать некоторые из Ваших фраз.\nСейчас режим вежливости %s\nЗа подробностями в /help"
			if GENTLE[update.Message.Chat.ID] {
				message = fmt.Sprintf(message, "включен")
			} else {
				message = fmt.Sprintf(message, "отключен")
			}
			sendMessage(w, update.Message.Chat.ID, message)
			return
		} else if strings.Contains(update.Message.Text, "/help") {
			sendMessage(w, update.Message.Chat.ID, "Для включения вежливого режима используйте команду /gentle\nДля отключения - /hardcore")
			return
		} else if strings.Contains(update.Message.Text, "/hardcore") {
			GENTLE[update.Message.Chat.ID] = false
			sendMessage(w, update.Message.Chat.ID, "Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle")
			return
		} else if strings.Contains(update.Message.Text, "/gentle") {
			GENTLE[update.Message.Chat.ID] = true
			sendMessage(w, update.Message.Chat.ID, "Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore")
			return
		} else {
			if _, ok := DELAY[update.Message.Chat.ID]; ok {
				DELAY[update.Message.Chat.ID] -= 1
			} else {
				DELAY[update.Message.Chat.ID] = rand.Intn(4)
			}
			if DELAY[update.Message.Chat.ID] == 0 {
				delete(DELAY, update.Message.Chat.ID)
				log.Infof(ctx, "[%v] %s", update.Message.Chat.ID, update.Message.Text)
				output := huify(update.Message.Text, GENTLE[update.Message.Chat.ID])
				if output != "" {
					sendMessage(w, update.Message.Chat.ID, output)
					return
				}
			}
		}
	})
}
