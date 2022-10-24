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
	settings.EnsureCache(request.ctx, request.cacheID)
	return request, nil
}

func (request *requestInfo) logWarn(err error) {
	log.Printf("[%v] %s", request.updateMessage.Chat.ID, err.Error())
}

func (request *requestInfo) answer(message, parseMode string) {
	sendMessage(request.writer, request.updateMessage.Chat.ID, message, nil, parseMode)
}

func (request *requestInfo) answerErrorWithLog(message string, err error, parseMode string) {
	request.logWarn(err)
	request.answer(message, parseMode)
}

func (request *requestInfo) isReplyNeeded() bool {
	return settings.cache[request.cacheID].Reply ||
		(request.updateMessage.ReplyTo != nil &&
			request.updateMessage.ReplyTo.From.Username != nil &&
			strings.Compare(*request.updateMessage.ReplyTo.From.Username, "xye_bot") == 0)
}

func (request *requestInfo) getReplyIDIfNeeded() *int64 {
	if request.isReplyNeeded() {
		return &request.updateMessage.ID
	}
	return nil
}

func (request *requestInfo) isStopped() bool {
	return !settings.cache[request.cacheID].Enabled
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
	case "/reply":
		command = &commandReply{request: request}
	case "/noreply":
		command = &commandNoReply{request: request}
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
		delayMap[request.updateMessage.Chat.ID] = rand.Intn(settings.cache[request.cacheID].Delay + 1)
	}
}

func (request *requestInfo) isAnswerNeeded() bool {
	return delayMap[request.updateMessage.Chat.ID] == 0
}

func (request *requestInfo) cleanDelay() {
	delete(delayMap, request.updateMessage.Chat.ID)
}

func (request *requestInfo) huify() string {
	return Huify(request.updateMessage.Text, settings.cache[request.cacheID].Gentle, rand.Intn(settings.cache[request.cacheID].WordsAmount)+1)
}
