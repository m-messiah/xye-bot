package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
)

func newRequest(w http.ResponseWriter, r *http.Request) (*requestInfo, error) {
	bytes, _ := ioutil.ReadAll(r.Body)
	var update Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		return nil, err
	}
	updateMessage := update.Message
	if updateMessage == nil {
		if update.EditedMessage == nil {
			return nil, errors.New("No message in update")
		}
		updateMessage = update.EditedMessage
	}
	request := &requestInfo{
		ctx:           r.Context(),
		updateMessage: updateMessage,
		writer:        w,
	}
	request.customDelayKey = datastore.NameKey("DatastoreDelay", strconv.FormatInt(updateMessage.Chat.ID, 10), nil)
	request.gentleKey = request.datastoreGetBool("Gentle")
	request.stoppedKey = request.datastoreGetBool("Stopped")
	request.wordsAmountKey = request.datastoreGetInt("WordsAmount")
	return request, nil
}

func (request *requestInfo) logWarn(err error) {
	log.Printf("[%v] %s", request.updateMessage.Chat.ID, err.Error())
}

func (request *requestInfo) answer(message string) {
	sendMessage(request.writer, request.updateMessage.Chat.ID, message, nil)
}

func (request *requestInfo) answerErrorWithLog(message string, err error) {
	request.logWarn(err)
	request.answer(message)
}

func (request *requestInfo) getReplyIDIfNeeded() *int64 {
	if request.updateMessage.ReplyTo != nil {
		if request.updateMessage.ReplyTo.From.Username != nil {
			if strings.Compare(*request.updateMessage.ReplyTo.From.Username, "xye_bot") == 0 {
				return &request.updateMessage.ID
			}
		}
	}
	return nil
}

func (request *requestInfo) datastoreGetBool(datastoreDBName string) *datastore.Key {
	var localCache map[int64]bool
	var resultStruct *DatastoreBool
	var defaultValue bool
	switch datastoreDBName {
	case "Gentle":
		localCache = gentleMap
		resultStruct = &request.gentleStruct
		defaultValue = true
	case "Stopped":
		localCache = stoppedMap
		resultStruct = &request.stoppedStruct
		defaultValue = false
	default:
		return nil
	}
	datastoreKey := datastore.NameKey(datastoreDBName, strconv.FormatInt(request.updateMessage.Chat.ID, 10), nil)
	if _, ok := localCache[request.updateMessage.Chat.ID]; !ok {
		if err := datastoreClient.Get(request.ctx, datastoreKey, resultStruct); err != nil {
			resultStruct.Value = defaultValue
			localCache[request.updateMessage.Chat.ID] = resultStruct.Value
			if _, err := datastoreClient.Put(request.ctx, datastoreKey, resultStruct); err != nil {
				log.Printf("[%v] %s %+v - %s", request.updateMessage.Chat.ID, datastoreKey, resultStruct, err.Error())
			}
		} else {
			localCache[request.updateMessage.Chat.ID] = resultStruct.Value
		}
	}
	return datastoreKey
}

func (request *requestInfo) datastoreGetInt(datastoreDBName string) *datastore.Key {
	var localCache map[int64]int
	var resultStruct *DatastoreInt
	var defaultValue int
	switch datastoreDBName {
	case "WordsAmount":
		localCache = wordsAmountMap
		resultStruct = &request.wordsAmountStruct
		defaultValue = 1
	default:
		return nil
	}
	datastoreKey := datastore.NameKey(datastoreDBName, strconv.FormatInt(request.updateMessage.Chat.ID, 10), nil)
	if _, ok := localCache[request.updateMessage.Chat.ID]; !ok {
		if err := datastoreClient.Get(request.ctx, datastoreKey, resultStruct); err != nil {
			resultStruct.Value = defaultValue
			localCache[request.updateMessage.Chat.ID] = resultStruct.Value
			if _, err := datastoreClient.Put(request.ctx, datastoreKey, resultStruct); err != nil {
				log.Printf("[%v] %s %+v - %s", request.updateMessage.Chat.ID, datastoreKey, resultStruct, err.Error())
			}
		} else {
			localCache[request.updateMessage.Chat.ID] = resultStruct.Value
		}
	}
	return datastoreKey
}

func (request *requestInfo) isStopped() bool {
	return stoppedMap[request.updateMessage.Chat.ID]
}

func (request *requestInfo) getStatusString() string {
	if request.isStopped() {
		return "остановлен"
	}
	return "включен"
}

func getCommandName(text string) string {
	commandName := ""
	if strings.Index(text, "/") == 0 {
		commandName = strings.Split(text, " ")[0]
		splittedCommand := strings.Split(commandName, "@")
		if len(splittedCommand) > 1 && splittedCommand[1] == "xye_bot" {
			commandName = splittedCommand[0]
		}
	}
	return commandName
}

func (request *requestInfo) getCommand() commandInterface {
	commandName := getCommandName(request.updateMessage.Text)
	var command commandInterface
	switch commandName {
	case "/start":
		command = &commandStart{request: request}
	case "/stop":
		command = &commandStop{request: request}
	case "/help":
		command = &commandHelp{request: request}
	case "/delay":
		command = &commandDelay{request: request}
	case "/hardcore":
		command = &commandHardcore{request: request}
	case "/gentle":
		command = &commandGentle{request: request}
	case "/amount":
		command = &commandAmount{request: request}
	default:
		command = &commandNotFound{request: request}
	}
	return command
}

func (request *requestInfo) parseCommand() error {
	return handleCommand(request.getCommand())
}

func (request *requestInfo) handleDelay() {
	if _, ok := delayMap[request.updateMessage.Chat.ID]; ok {
		delayMap[request.updateMessage.Chat.ID]--
	} else {
		if currentDelay, ok := customDelayMap[request.updateMessage.Chat.ID]; ok {
			delayMap[request.updateMessage.Chat.ID] = rand.Intn(currentDelay + 1)
		} else {
			if err := datastoreClient.Get(request.ctx, request.customDelayKey, &request.customDelay); err != nil {
				log.Printf("[%v] %s - %s", request.updateMessage.Chat.ID, request.customDelayKey, err.Error())
				request.customDelay.Delay = defaultDelay
				customDelayMap[request.updateMessage.Chat.ID] = defaultDelay
				if _, err := datastoreClient.Put(request.ctx, request.customDelayKey, &request.customDelay); err != nil {
					log.Printf("[%v] %s - %s", request.updateMessage.Chat.ID, request.customDelayKey, err.Error())
				}
			} else {
				customDelayMap[request.updateMessage.Chat.ID] = request.customDelay.Delay
				delayMap[request.updateMessage.Chat.ID] = rand.Intn(request.customDelay.Delay + 1)
			}
		}
	}
}

func (request *requestInfo) isAnswerNeeded(replyID *int64) bool {
	return delayMap[request.updateMessage.Chat.ID] == 0 || replyID != nil
}

func (request *requestInfo) cleanDelay() {
	delete(delayMap, request.updateMessage.Chat.ID)
}

func (request *requestInfo) huify() string {
	return Huify(request.updateMessage.Text, gentleMap[request.updateMessage.Chat.ID], rand.Intn(wordsAmountMap[request.updateMessage.Chat.ID])+1)
}
