package xyebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func NewRequest(r *http.Request) (*Request, error) {
	self := &Request{}
	self.ctx = appengine.NewContext(r)
	bytes, _ := ioutil.ReadAll(r.Body)
	var update Update
	json.Unmarshal(bytes, &update)
	updateMessage := update.Message
	if updateMessage == nil {
		if update.EditedMessage == nil {
			return nil, errors.New("No message in update")
		}
		updateMessage = update.EditedMessage
	}
	self.updateMessage = updateMessage
	self.customDelayKey = datastore.NewKey(self.ctx, "DatastoreDelay", "", updateMessage.Chat.ID, nil)
	self.gentleKey = self.DatastoreGetBool("Gentle")
	self.stoppedKey = self.DatastoreGetBool("Stopped")
	self.wordsAmountKey = self.DatastoreGetInt("WordsAmount")
	return self, nil
}

func (self *Request) GetReplyIDIfNeeded() *int64 {
	if self.updateMessage.ReplyTo != nil {
		if self.updateMessage.ReplyTo.From.Username != nil {
			if strings.Compare(*self.updateMessage.ReplyTo.From.Username, "xye_bot") == 0 {
				return &self.updateMessage.ID
			}
		}
	}
	return nil
}

func (self *Request) DatastoreGetBool(datastoreDBName string) *datastore.Key {
	var localCache map[int64]bool
	var resultStruct *DatastoreBool
	var defaultValue bool
	switch datastoreDBName {
	case "Gentle":
		localCache = Gentle
		resultStruct = &self.gentleStruct
		defaultValue = true
	case "Stopped":
		localCache = Stopped
		resultStruct = &self.stoppedStruct
		defaultValue = false
	default:
		return nil
	}
	datastoreKey := datastore.NewKey(self.ctx, datastoreDBName, "", self.updateMessage.Chat.ID, nil)
	if _, ok := localCache[self.updateMessage.Chat.ID]; !ok {
		if err := datastore.Get(self.ctx, datastoreKey, &resultStruct); err != nil {
			resultStruct.Value = defaultValue
			localCache[self.updateMessage.Chat.ID] = resultStruct.Value
			if _, err := datastore.Put(self.ctx, datastoreKey, &resultStruct); err != nil {
				log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
			}
		} else {
			localCache[self.updateMessage.Chat.ID] = resultStruct.Value
		}
	}
	return datastoreKey
}

func (self *Request) DatastoreGetInt(datastoreDBName string) *datastore.Key {
	var localCache map[int64]int
	var resultStruct *DatastoreInt
	var defaultValue int
	switch datastoreDBName {
	case "WordsAmount":
		localCache = WordsAmount
		resultStruct = &self.wordsAmountStruct
		defaultValue = 1
	default:
		return nil
	}
	datastoreKey := datastore.NewKey(self.ctx, datastoreDBName, "", self.updateMessage.Chat.ID, nil)
	if _, ok := localCache[self.updateMessage.Chat.ID]; !ok {
		if err := datastore.Get(self.ctx, datastoreKey, &resultStruct); err != nil {
			resultStruct.Value = defaultValue
			localCache[self.updateMessage.Chat.ID] = resultStruct.Value
			if _, err := datastore.Put(self.ctx, datastoreKey, &resultStruct); err != nil {
				log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
			}
		} else {
			localCache[self.updateMessage.Chat.ID] = resultStruct.Value
		}
	}
	return datastoreKey
}

func (self *Request) IsStopped() bool {
	return Stopped[self.updateMessage.Chat.ID]
}

func (self *Request) ParseCommand(w http.ResponseWriter) error {
	if IsCommand(self.updateMessage.Text, "/start") {
		message := "Привет! Я бот-хуебот.\n" +
			"Я буду хуифицировать некоторые из Ваших фраз.\n" +
			"Сейчас режим вежливости %s\n" +
			"За подробностями в /help"
		Stopped[self.updateMessage.Chat.ID] = false
		self.stoppedStruct.Value = false
		if _, err := datastore.Put(self.ctx, self.stoppedKey, &self.stoppedStruct); err != nil {
			log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
		}
		if Gentle[self.updateMessage.Chat.ID] {
			message = fmt.Sprintf(message, "включен")
		} else {
			message = fmt.Sprintf(message, "отключен")
		}
		SendMessage(w, self.updateMessage.Chat.ID, message, nil)
		return nil
	}

	if IsCommand(self.updateMessage.Text, "/stop") {
		Stopped[self.updateMessage.Chat.ID] = true
		self.stoppedStruct.Value = true
		if _, err := datastore.Put(self.ctx, self.stoppedKey, &self.stoppedStruct); err != nil {
			log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
		}
		SendMessage(w, self.updateMessage.Chat.ID, "Выключаюсь", nil)
		return nil
	}

	if IsCommand(self.updateMessage.Text, "/help") {
		SendMessage(w, self.updateMessage.Chat.ID,
			"Вежливый режим:\n"+
				"  Для включения используйте команду /gentle\n"+
				"  Для отключения - /hardcore\n"+
				"Частота ответов: /delay N, где N - любое любое натуральное число\n"+
				"Число хуифицируемых слов: /amount N, где N - от 1 до 10\n"+
				"Для остановки используйте /stop", nil)
		return nil
	}
	if IsCommand(self.updateMessage.Text, "/delay") {
		command := strings.Fields(self.updateMessage.Text)
		if len(command) < 2 {
			currentDelayMessage := "Сейчас я пропускаю случайное число сообщений от 0 до "
			if currentDelay, ok := CustomDelay[self.updateMessage.Chat.ID]; ok {
				currentDelayMessage += strconv.Itoa(currentDelay)
			} else {
				currentDelayMessage += "4"
			}
			SendMessage(w, self.updateMessage.Chat.ID, currentDelayMessage, nil)
			return nil
		}
		commandArg := command[len(command)-1]
		tryDelay, err := strconv.Atoi(commandArg)
		if err != nil || tryDelay < 1 || tryDelay > 1000000 {
			SendMessage(w, self.updateMessage.Chat.ID, "Неправильный аргумент, отправьте `/delay N`, где N любое натуральное число меньше 1000000", nil)
			return nil
		}
		self.customDelay.Delay = tryDelay
		if _, err := datastore.Put(self.ctx, self.customDelayKey, &self.customDelay); err != nil {
			log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
			SendMessage(w, self.updateMessage.Chat.ID, "Не удалось сохранить, отправьте еще раз `/delay N`, где N любое натуральное число меньше 1000000", nil)
			return nil
		}
		CustomDelay[self.updateMessage.Chat.ID] = self.customDelay.Delay
		SendMessage(w, self.updateMessage.Chat.ID, "Я буду пропускать случайное число сообщений от 0 до "+commandArg, nil)
		delete(Delay, self.updateMessage.Chat.ID)
		return nil
	}
	if IsCommand(self.updateMessage.Text, "/hardcore") {
		Gentle[self.updateMessage.Chat.ID] = false
		self.gentleStruct.Value = false
		if _, err := datastore.Put(self.ctx, self.gentleKey, &self.gentleStruct); err != nil {
			log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
		}
		SendMessage(w, self.updateMessage.Chat.ID, "Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle", nil)
		return nil
	}
	if IsCommand(self.updateMessage.Text, "/gentle") {
		Gentle[self.updateMessage.Chat.ID] = true
		self.gentleStruct.Value = true
		if _, err := datastore.Put(self.ctx, self.gentleKey, &self.gentleStruct); err != nil {
			log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
		}
		SendMessage(w, self.updateMessage.Chat.ID, "Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore", nil)
		return nil
	}
	if IsCommand(self.updateMessage.Text, "/amount") {
		command := strings.Fields(self.updateMessage.Text)
		if len(command) < 2 {
			currentWordsAmount := 1
			if currentAmount, ok := WordsAmount[self.updateMessage.Chat.ID]; ok {
				currentWordsAmount = currentAmount
			}
			SendMessage(w, self.updateMessage.Chat.ID, "Сейчас я хуифицирую случайное число слов от 1 до "+strconv.Itoa(currentWordsAmount), nil)
			return nil
		}
		commandArg := command[len(command)-1]
		tryWordsAmount, err := strconv.Atoi(commandArg)
		if err != nil || tryWordsAmount < 1 || tryWordsAmount > 10 {
			SendMessage(w, self.updateMessage.Chat.ID, "Неправильный аргумент, отправьте `/amount N`, где N любое натуральное число не больше 10", nil)
			return nil
		}
		self.wordsAmountStruct.Value = tryWordsAmount
		if _, err := datastore.Put(self.ctx, self.wordsAmountKey, &self.wordsAmountStruct); err != nil {
			log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
			SendMessage(w, self.updateMessage.Chat.ID, "Не удалось сохранить, отправьте еще раз `/amount N`, где N любое натуральное число не больше 10", nil)
			return nil
		}
		WordsAmount[self.updateMessage.Chat.ID] = self.wordsAmountStruct.Value
		SendMessage(w, self.updateMessage.Chat.ID, "Я буду хуифицировать случайное число слов от 1 до "+strconv.Itoa(self.wordsAmountStruct.Value), nil)
		return nil
	}
	return errors.New("Команда не найдена")
}

func (self *Request) HandleDelay() {
	if _, ok := Delay[self.updateMessage.Chat.ID]; ok {
		Delay[self.updateMessage.Chat.ID]--
	} else {
		if currentDelay, ok := CustomDelay[self.updateMessage.Chat.ID]; ok {
			Delay[self.updateMessage.Chat.ID] = rand.Intn(currentDelay + 1)
		} else {
			if err := datastore.Get(self.ctx, self.customDelayKey, &self.customDelay); err != nil {
				log.Infof(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
				self.customDelay.Delay = DEFAULT_DELAY
				CustomDelay[self.updateMessage.Chat.ID] = DEFAULT_DELAY
				if _, err := datastore.Put(self.ctx, self.customDelayKey, &self.customDelay); err != nil {
					log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
				}
			} else {
				CustomDelay[self.updateMessage.Chat.ID] = self.customDelay.Delay
				Delay[self.updateMessage.Chat.ID] = rand.Intn(self.customDelay.Delay + 1)
			}
		}
	}
}

func (self *Request) IsAnswerNeeded(replyID *int64) bool {
	return Delay[self.updateMessage.Chat.ID] == 0 || replyID != nil
}

func (self *Request) CleanDelay() {
	delete(Delay, self.updateMessage.Chat.ID)
}

func (self *Request) Huify() string {
	return Huify(self.updateMessage.Text, Gentle[self.updateMessage.Chat.ID], rand.Intn(WordsAmount[self.updateMessage.Chat.ID])+1)
}
