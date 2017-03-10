package xyebot

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	CUSTOM_DELAY := make(map[int64]int)
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		ctx := appengine.NewContext(r)
		var custom_delay DatastoreDelay
		var gentle_struct DatastoreGentle
		var update Update
		json.Unmarshal(bytes, &update)
		if update.Message == nil {
			return
		}
		custom_delay_key := datastore.NewKey(ctx, "DatastoreDelay", "", update.Message.Chat.ID, nil)
		gentle_key := datastore.NewKey(ctx, "Gentle", "", update.Message.Chat.ID, nil)
		if _, ok := GENTLE[update.Message.Chat.ID]; !ok {
			if err := datastore.Get(ctx, gentle_key, &gentle_struct); err != nil {
				GENTLE[update.Message.Chat.ID] = true
				gentle_struct.Gentle = true
				if _, err := datastore.Put(ctx, gentle_key, &gentle_struct); err != nil {
					log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
				}
			} else {
				GENTLE[update.Message.Chat.ID] = gentle_struct.Gentle
			}

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
					"Частота ответов: /delay N, где N - любое любое натуральное число")
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
			if err != nil || try_delay < 1 {
				sendMessage(w, update.Message.Chat.ID, "Неправильный аргумент, отправьте `/delay N`, где N любое натуральное число")
				return
			}
			custom_delay.Delay = try_delay
			if _, err := datastore.Put(ctx, custom_delay_key, &custom_delay); err != nil {
				log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
				sendMessage(w, update.Message.Chat.ID, "Не удалось сохранить, отправьте еще раз `/delay N`, где N любое натуральное число")
				return
			}
			CUSTOM_DELAY[update.Message.Chat.ID] = custom_delay.Delay
			sendMessage(w, update.Message.Chat.ID, "Я буду пропускать случайное число сообщений от 0 до "+command_arg)
			return
		} else if strings.Contains(update.Message.Text, "/hardcore") {
			GENTLE[update.Message.Chat.ID] = false
			gentle_struct.Gentle = false
			if _, err := datastore.Put(ctx, gentle_key, &gentle_struct); err != nil {
				log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
			}
			sendMessage(w, update.Message.Chat.ID, "Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle")
			return
		} else if strings.Contains(update.Message.Text, "/gentle") {
			GENTLE[update.Message.Chat.ID] = true
			gentle_struct.Gentle = true
			if _, err := datastore.Put(ctx, gentle_key, &gentle_struct); err != nil {
				log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
			}
			sendMessage(w, update.Message.Chat.ID, "Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore")
			return
		} else {
			if _, ok := DELAY[update.Message.Chat.ID]; ok {
				DELAY[update.Message.Chat.ID] -= 1
			} else {
				if current_delay, ok := CUSTOM_DELAY[update.Message.Chat.ID]; ok {
					DELAY[update.Message.Chat.ID] = rand.Intn(current_delay + 1)
				} else {
					if err := datastore.Get(ctx, custom_delay_key, &custom_delay); err != nil {
						custom_delay.Delay = 4
						CUSTOM_DELAY[update.Message.Chat.ID] = 4
						if _, err := datastore.Put(ctx, custom_delay_key, &custom_delay); err != nil {
							log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
						}
					} else {
						CUSTOM_DELAY[update.Message.Chat.ID] = custom_delay.Delay
						DELAY[update.Message.Chat.ID] = rand.Intn(custom_delay.Delay + 1)
					}
				}
			}
			if DELAY[update.Message.Chat.ID] == 0 {
				delete(DELAY, update.Message.Chat.ID)
				// log.Infof(ctx, "[%v] %s", update.Message.Chat.ID, update.Message.Text)
				output := huify(update.Message.Text, GENTLE[update.Message.Chat.ID])
				if output != "" {
					sendMessage(w, update.Message.Chat.ID, output)
					return
				}
			}
		}
	})
}
