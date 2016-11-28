package xyebot

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
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
	CUSTOM_DELAY := make(map[int64]int)
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
			sendMessage(w, update.Message.Chat.ID,
				"Вежливый режим:\n"+
					"  Для включения используйте команду /gentle\n"+
					"  Для отключения - /hardcore\n"+
					"Частота ответов: /delay N, где N - любое любое целое положительное число")
			return
		} else if strings.Contains(update.Message.Text, "/delay") {
			command := strings.Fields(update.Message.Text)
			if len(command) < 2 {
				current_delay_message := "Сейчас я пропускаю случайное число сообщений от 0 до "
				if current_delay, ok := CUSTOM_DELAY[update.Message.Chat.ID]; ok {
					current_delay_message += strconv.Itoa(current_delay)
				} else {
					current_delay_message += "4"
				}
				sendMessage(w, update.Message.Chat.ID, current_delay_message)
				return
			}
			command_arg := command[len(command)-1]
			try_delay, err := strconv.Atoi(command_arg)
			if err != nil || try_delay < 0 {
				sendMessage(w, update.Message.Chat.ID, "Неправильный аргумент, отправьте `/delay N`, где N любое целое положительное число")
				return
			}
			CUSTOM_DELAY[update.Message.Chat.ID] = try_delay
			sendMessage(w, update.Message.Chat.ID, "Я буду пропускать случайное число сообщений от 0 до "+command_arg)
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
				if custom_delay, ok := CUSTOM_DELAY[update.Message.Chat.ID]; ok {
					DELAY[update.Message.Chat.ID] = rand.Intn(custom_delay)
				} else {
					DELAY[update.Message.Chat.ID] = rand.Intn(4)
				}
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
