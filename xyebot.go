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
	Stopped := make(map[int64]bool)
	CustomDelay := make(map[int64]int)
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		ctx := appengine.NewContext(r)
		var customDelay DatastoreDelay
		var gentleStruct DatastoreBool
		var stoppedStruct DatastoreBool
		var update Update
		json.Unmarshal(bytes, &update)
		updateMessage := update.Message
		if updateMessage == nil {
			if update.EditedMessage == nil {
				return
			}
			updateMessage = update.EditedMessage
		}
		customDelayKey := datastore.NewKey(ctx, "DatastoreDelay", "", updateMessage.Chat.ID, nil)
		gentleKey := datastore.NewKey(ctx, "Gentle", "", updateMessage.Chat.ID, nil)
		if _, ok := Gentle[updateMessage.Chat.ID]; !ok {
			if err := datastore.Get(ctx, gentleKey, &gentleStruct); err != nil {
				Gentle[updateMessage.Chat.ID] = true
				gentleStruct.Value = true
				if _, err := datastore.Put(ctx, gentleKey, &gentleStruct); err != nil {
					log.Warningf(ctx, "[%v] %s", updateMessage.Chat.ID, err.Error())
				}
			} else {
				Gentle[updateMessage.Chat.ID] = gentleStruct.Value
			}
		}

		stoppedKey := datastore.NewKey(ctx, "Stopped", "", updateMessage.Chat.ID, nil)
		if _, ok := Stopped[updateMessage.Chat.ID]; !ok {
			if err := datastore.Get(ctx, stoppedKey, &stoppedStruct); err != nil {
				Stopped[updateMessage.Chat.ID] = false
				stoppedStruct.Value = false
				if _, err := datastore.Put(ctx, stoppedKey, &stoppedStruct); err != nil {
					log.Warningf(ctx, "[%v] %s", updateMessage.Chat.ID, err.Error())
				}
			} else {
				Stopped[updateMessage.Chat.ID] = stoppedStruct.Value
			}

		}

		if isCommand(updateMessage.Text, "/start") {
			message := "Привет! Я бот-хуебот.\nЯ буду хуифицировать некоторые из Ваших фраз.\nСейчас режим вежливости %s\nЗа подробностями в /help"
			Stopped[updateMessage.Chat.ID] = false
			stoppedStruct.Value = false
			if _, err := datastore.Put(ctx, stoppedKey, &stoppedStruct); err != nil {
				log.Warningf(ctx, "[%v] %s", updateMessage.Chat.ID, err.Error())
			}
			if Gentle[updateMessage.Chat.ID] {
				message = fmt.Sprintf(message, "включен")
			} else {
				message = fmt.Sprintf(message, "отключен")
			}
			sendMessage(w, updateMessage.Chat.ID, message)
			return
		}

		if isCommand(updateMessage.Text, "/stop") {
			Stopped[updateMessage.Chat.ID] = true
			stoppedStruct.Value = true
			if _, err := datastore.Put(ctx, stoppedKey, &stoppedStruct); err != nil {
				log.Warningf(ctx, "[%v] %s", updateMessage.Chat.ID, err.Error())
			}
			sendMessage(w, updateMessage.Chat.ID, "Выключаюсь")
			return
		}

		if isCommand(updateMessage.Text, "/help") {
			sendMessage(w, updateMessage.Chat.ID,
				"Вежливый режим:\n"+
					"  Для включения используйте команду /gentle\n"+
					"  Для отключения - /hardcore\n"+
					"Частота ответов: /delay N, где N - любое любое натуральное число\n"+
					"Для остановки используйте /stop")
			return
		}
		if isCommand(updateMessage.Text, "/delay") {
			command := strings.Fields(updateMessage.Text)
			if len(command) < 2 {
				currentDelayMessage := "Сейчас я пропускаю случайное число сообщений от 0 до "
				if currentDelay, ok := CustomDelay[updateMessage.Chat.ID]; ok {
					currentDelayMessage += strconv.Itoa(currentDelay)
				} else {
					currentDelayMessage += "4"
				}
				sendMessage(w, updateMessage.Chat.ID, currentDelayMessage)
				return
			}
			commandArg := command[len(command)-1]
			tryDelay, err := strconv.Atoi(commandArg)
			if err != nil || tryDelay < 1 || tryDelay > 1000000 {
				sendMessage(w, updateMessage.Chat.ID, "Неправильный аргумент, отправьте `/delay N`, где N любое натуральное число меньше 1000000")
				return
			}
			customDelay.Delay = tryDelay
			if _, err := datastore.Put(ctx, customDelayKey, &customDelay); err != nil {
				log.Warningf(ctx, "[%v] %s", updateMessage.Chat.ID, err.Error())
				sendMessage(w, updateMessage.Chat.ID, "Не удалось сохранить, отправьте еще раз `/delay N`, где N любое натуральное число меньше 1000000")
				return
			}
			CustomDelay[updateMessage.Chat.ID] = customDelay.Delay
			sendMessage(w, updateMessage.Chat.ID, "Я буду пропускать случайное число сообщений от 0 до "+commandArg)
			return
		}
		if isCommand(updateMessage.Text, "/hardcore") {
			Gentle[updateMessage.Chat.ID] = false
			gentleStruct.Value = false
			if _, err := datastore.Put(ctx, gentleKey, &gentleStruct); err != nil {
				log.Warningf(ctx, "[%v] %s", updateMessage.Chat.ID, err.Error())
			}
			sendMessage(w, updateMessage.Chat.ID, "Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle")
			return
		}
		if isCommand(updateMessage.Text, "/gentle") {
			Gentle[updateMessage.Chat.ID] = true
			gentleStruct.Value = true
			if _, err := datastore.Put(ctx, gentleKey, &gentleStruct); err != nil {
				log.Warningf(ctx, "[%v] %s", updateMessage.Chat.ID, err.Error())
			}
			sendMessage(w, updateMessage.Chat.ID, "Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore")
			return
		}

		if Stopped[updateMessage.Chat.ID] {
			return
		}

		if _, ok := Delay[updateMessage.Chat.ID]; ok {
			Delay[updateMessage.Chat.ID]--
		} else {
			if currentDelay, ok := CustomDelay[updateMessage.Chat.ID]; ok {
				Delay[updateMessage.Chat.ID] = rand.Intn(currentDelay + 1)
			} else {
				if err := datastore.Get(ctx, customDelayKey, &customDelay); err != nil {
					customDelay.Delay = 4
					CustomDelay[updateMessage.Chat.ID] = 4
					if _, err := datastore.Put(ctx, customDelayKey, &customDelay); err != nil {
						log.Warningf(ctx, "[%v] %s", updateMessage.Chat.ID, err.Error())
					}
				} else {
					CustomDelay[updateMessage.Chat.ID] = customDelay.Delay
					Delay[updateMessage.Chat.ID] = rand.Intn(customDelay.Delay + 1)
				}
			}
		}
		if Delay[updateMessage.Chat.ID] == 0 {
			delete(Delay, updateMessage.Chat.ID)
			// log.Infof(ctx, "[%v] %s", updateMessage.Chat.ID, updateMessage.Text)
			output := huify(updateMessage.Text, Gentle[updateMessage.Chat.ID])
			if output != "" {
				sendMessage(w, updateMessage.Chat.ID, output)
				return
			}
		}
	})
}
