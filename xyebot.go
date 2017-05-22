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

func sendMessage(w http.ResponseWriter, chatID int64, text string) {
	msg := Response{Chatid: chatID, Text: text, Method: "sendMessage"}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

func isCommand(text, command string) bool {
	if strings.Index(text, command) == 0 {
		if strings.Contains(text, "@xye_bot") || !strings.Contains(text, "@") {
			return true
		}
	}
	return false
}

func init() {
	Delay := make(map[int64]int)
	Gentle := make(map[int64]bool)
	CustomDelay := make(map[int64]int)
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		ctx := appengine.NewContext(r)
		var customDelay DatastoreDelay
		var gentleStruct DatastoreGentle
		var update Update
		json.Unmarshal(bytes, &update)
		if update.Message == nil {
			return
		}
		customDelayKey := datastore.NewKey(ctx, "DatastoreDelay", "", update.Message.Chat.ID, nil)
		gentleKey := datastore.NewKey(ctx, "Gentle", "", update.Message.Chat.ID, nil)
		if _, ok := Gentle[update.Message.Chat.ID]; !ok {
			if err := datastore.Get(ctx, gentleKey, &gentleStruct); err != nil {
				Gentle[update.Message.Chat.ID] = true
				gentleStruct.Gentle = true
				if _, err := datastore.Put(ctx, gentleKey, &gentleStruct); err != nil {
					log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
				}
			} else {
				Gentle[update.Message.Chat.ID] = gentleStruct.Gentle
			}

		}

		if isCommand(update.Message.Text, "/start") {
			message := "Привет! Я бот-хуебот.\nЯ буду хуифицировать некоторые из Ваших фраз.\nСейчас режим вежливости %s\nЗа подробностями в /help"
			if Gentle[update.Message.Chat.ID] {
				message = fmt.Sprintf(message, "включен")
			} else {
				message = fmt.Sprintf(message, "отключен")
			}
			sendMessage(w, update.Message.Chat.ID, message)
			return
		}

		if isCommand(update.Message.Text, "/help") {
			sendMessage(w, update.Message.Chat.ID,
				"Вежливый режим:\n"+
					"  Для включения используйте команду /gentle\n"+
					"  Для отключения - /hardcore\n"+
					"Частота ответов: /delay N, где N - любое любое натуральное число")
			return
		}
		if isCommand(update.Message.Text, "/delay") {
			command := strings.Fields(update.Message.Text)
			if len(command) < 2 {
				currentDelayMessage := "Сейчас я пропускаю случайное число сообщений от 0 до "
				if currentDelay, ok := CustomDelay[update.Message.Chat.ID]; ok {
					currentDelayMessage += strconv.Itoa(currentDelay)
				} else {
					currentDelayMessage += "4"
				}
				sendMessage(w, update.Message.Chat.ID, currentDelayMessage)
				return
			}
			commandArg := command[len(command)-1]
			tryDelay, err := strconv.Atoi(commandArg)
			if err != nil || tryDelay < 1 {
				sendMessage(w, update.Message.Chat.ID, "Неправильный аргумент, отправьте `/delay N`, где N любое натуральное число")
				return
			}
			customDelay.Delay = tryDelay
			if _, err := datastore.Put(ctx, customDelayKey, &customDelay); err != nil {
				log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
				sendMessage(w, update.Message.Chat.ID, "Не удалось сохранить, отправьте еще раз `/delay N`, где N любое натуральное число")
				return
			}
			CustomDelay[update.Message.Chat.ID] = customDelay.Delay
			sendMessage(w, update.Message.Chat.ID, "Я буду пропускать случайное число сообщений от 0 до "+commandArg)
			return
		}
		if isCommand(update.Message.Text, "/hardcore") {
			Gentle[update.Message.Chat.ID] = false
			gentleStruct.Gentle = false
			if _, err := datastore.Put(ctx, gentleKey, &gentleStruct); err != nil {
				log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
			}
			sendMessage(w, update.Message.Chat.ID, "Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle")
			return
		}
		if isCommand(update.Message.Text, "/gentle") {
			Gentle[update.Message.Chat.ID] = true
			gentleStruct.Gentle = true
			if _, err := datastore.Put(ctx, gentleKey, &gentleStruct); err != nil {
				log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
			}
			sendMessage(w, update.Message.Chat.ID, "Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore")
			return
		}

		if _, ok := Delay[update.Message.Chat.ID]; ok {
			Delay[update.Message.Chat.ID]--
		} else {
			if currentDelay, ok := CustomDelay[update.Message.Chat.ID]; ok {
				Delay[update.Message.Chat.ID] = rand.Intn(currentDelay + 1)
			} else {
				if err := datastore.Get(ctx, customDelayKey, &customDelay); err != nil {
					customDelay.Delay = 4
					CustomDelay[update.Message.Chat.ID] = 4
					if _, err := datastore.Put(ctx, customDelayKey, &customDelay); err != nil {
						log.Warningf(ctx, "[%v] %s", update.Message.Chat.ID, err.Error())
					}
				} else {
					CustomDelay[update.Message.Chat.ID] = customDelay.Delay
					Delay[update.Message.Chat.ID] = rand.Intn(customDelay.Delay + 1)
				}
			}
		}
		if Delay[update.Message.Chat.ID] == 0 {
			delete(Delay, update.Message.Chat.ID)
			// log.Infof(ctx, "[%v] %s", update.Message.Chat.ID, update.Message.Text)
			output := huify(update.Message.Text, Gentle[update.Message.Chat.ID])
			if output != "" {
				sendMessage(w, update.Message.Chat.ID, output)
				return
			}
		}
	})
}
