package xyebot

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
    "fmt"
    "os"
)

func sendMessage(w http.ResponseWriter, chat_id int64, text string) {
    msg := Response{Chatid: chat_id, Text: text, Method: "sendMessage"}
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(msg)
}

func init() {
	DELAY := make(map[int64]int)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)

		var update Update
		json.Unmarshal(bytes, &update)
        if update.Message == nil {
			return
		}

		if strings.Contains(update.Message.Text, "/start") || strings.Contains(update.Message.Text, "/help") {
			sendMessage(w, update.Message.Chat.ID, "Привет! Я бот-хуебот.\nЯ буду хуифицировать некоторые из твоих фраз")
			return
		} else {
			if _, ok := DELAY[update.Message.Chat.ID]; ok {
				DELAY[update.Message.Chat.ID] -= 1
			} else {
				DELAY[update.Message.Chat.ID] = rand.Intn(4)
			}
			if DELAY[update.Message.Chat.ID] == 0 {
				delete(DELAY, update.Message.Chat.ID)
                fmt.Fprintf(os.Stderr, "[%v] %s\n", update.Message.Chat.ID, string(bytes))
				output := huify(update.Message.Text)
				if output != "" {
					sendMessage(w, update.Message.Chat.ID, output)
					return
				}
			}
		}
	})
}
