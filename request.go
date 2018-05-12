package xyebot

import (
	"encoding/json"
	"errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

func NewRequest(w http.ResponseWriter, r *http.Request) (*Request, error) {
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
	self := &Request{
		ctx:           appengine.NewContext(r),
		updateMessage: updateMessage,
		writer:        w,
	}
	self.customDelayKey = datastore.NewKey(self.ctx, "DatastoreDelay", "", updateMessage.Chat.ID, nil)
	self.gentleKey = self.DatastoreGetBool("Gentle")
	self.stoppedKey = self.DatastoreGetBool("Stopped")
	self.wordsAmountKey = self.DatastoreGetInt("WordsAmount")
	return self, nil
}

func (self *Request) LogWarn(err error) {
	log.Warningf(self.ctx, "[%v] %s", self.updateMessage.Chat.ID, err.Error())
}

func (self *Request) Answer(message string) {
	SendMessage(self.writer, self.updateMessage.Chat.ID, message, nil)
}

func (self *Request) AnswerErrorWithLog(message string, err error) {
	self.LogWarn(err)
	self.Answer(message)
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

func GetCommandName(text string) string {
	command_name := ""
	if strings.Index(text, "/") == 0 {
		command_name = strings.Split(text, " ")[0]
		splitted_command := strings.Split(command_name, "@")
		if len(splitted_command) > 1 && splitted_command[1] == "xye_bot" {
			command_name = splitted_command[0]
		}
	}
	return command_name
}

func (self *Request) GetCommand() CommandIF {
	command_name := GetCommandName(self.updateMessage.Text)
	var command CommandIF
	switch command_name {
	case "/start":
		command = &CommandStart{request: self}
	case "/stop":
		command = &CommandStop{request: self}
	case "/help":
		command = &CommandHelp{request: self}
	case "/delay":
		command = &CommandDelay{request: self}
	case "/hardcore":
		command = &CommandHardcore{request: self}
	case "/gentle":
		command = &CommandGentle{request: self}
	case "/amount":
		command = &CommandAmount{request: self}
	default:
		command = &CommandNotFound{request: self}
	}
	return command
}

func (self *Request) ParseCommand(w http.ResponseWriter) error {
	return handleCommand(self.GetCommand())
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
