package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func newRequest(w http.ResponseWriter, r *http.Request) (*requestInfo, error) {
	bytes, _ := io.ReadAll(r.Body)
	var update Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		return nil, err
	}
	updateMessage := update.Message
	if updateMessage == nil {
		if update.EditedMessage == nil {
			return nil, errors.New("no message in update")
		}
		updateMessage = update.EditedMessage
	}
	request := &requestInfo{
		ctx:           r.Context(),
		updateMessage: updateMessage,
		writer:        w,
		cacheID:       strconv.FormatInt(updateMessage.Chat.ID, 10),
	}
	return request, nil
}

func (request *requestInfo) logWarn(err error) {
	log.Printf("[%v] %s", request.updateMessage.Chat.ID, err.Error())
}

func (request *requestInfo) answer(message, parseMode string) {
	sendMessage(request.writer, request.updateMessage.Chat.ID, message, &request.updateMessage.ID, parseMode)
}

func getCommandName(text string) string {
	commandName := ""
	if strings.Index(text, "/") == 0 {
		commandName = strings.Split(text, " ")[0]
		commandParts := strings.Split(commandName, "@")
		if len(commandParts) > 1 && commandParts[1] == BotName {
			commandName = commandParts[0]
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
	default:
		command = &commandNotFound{request: request}
	}
	return command
}

func (request *requestInfo) handleCommand() error {
	return handleCommand(request.getCommand())
}

func (request *requestInfo) handleDelay() {
	if _, ok := delayMap[request.updateMessage.Chat.ID]; ok {
		delayMap[request.updateMessage.Chat.ID]--
	} else {
		delayMap[request.updateMessage.Chat.ID] = rand.Intn(DelayLimit)
	}
}

func (request *requestInfo) isAnswerNeeded() bool {
	return delayMap[request.updateMessage.Chat.ID] == 0
}

func (request *requestInfo) cleanDelay() {
	delete(delayMap, request.updateMessage.Chat.ID)
}

func (request *requestInfo) Modify() string {
	return Huify(request.updateMessage.Text, WordsAmount)
}
