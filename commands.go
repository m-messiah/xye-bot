package main

import (
	"cloud.google.com/go/datastore"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type commandStart botCommand

func (commandRequest *commandStart) Handle() error {
	message := "Привет! Я бот-хуебот.\n" +
		"Я буду хуифицировать некоторые из ваших фраз.\n" +
		"Сейчас режим вежливости %s\n" +
		"За подробностями в /help."
	switchDatastoreBool(commandRequest.request, "Stopped", false)
	if gentleMap[commandRequest.request.updateMessage.Chat.ID] {
		message = fmt.Sprintf(message, "включен")
	} else {
		message = fmt.Sprintf(message, "отключен")
	}
	commandRequest.request.answer(message)
	return nil
}

type commandStop botCommand

func (commandRequest *commandStop) Handle() error {
	switchDatastoreBool(commandRequest.request, "Stopped", true)
	commandRequest.request.answer("Выключаюсь")
	return nil
}

type commandHelp botCommand

func (commandRequest *commandHelp) Handle() error {
	commandRequest.request.answer(
		"Вежливый режим:\n" +
			"  Для включения используйте команду /gentle\n" +
			"  Для отключения - /hardcore\n" +
			"Частота ответов: /delay N, где N - любое любое натуральное число\n" +
			"Число хуифицируемых слов: /amount N, где N - от 1 до 10\n" +
			"Для остановки используйте /stop\n\n" + 
			"По вопросам: https://github.com/m-messiah/xye-bot/issues")
	return nil
}

type commandDelay botCommand

func (commandRequest *commandDelay) Handle() error {
	command := strings.Fields(commandRequest.request.updateMessage.Text)
	if len(command) < 2 {
		currentDelayMessage := "Сейчас я пропускаю случайное число сообщений от 0 до "
		if currentDelay, ok := customDelayMap[commandRequest.request.updateMessage.Chat.ID]; ok {
			currentDelayMessage += strconv.Itoa(currentDelay)
		} else {
			currentDelayMessage += "4"
		}
		commandRequest.request.answer(currentDelayMessage)
		return nil
	}
	commandArg := command[len(command)-1]
	tryDelay, err := strconv.Atoi(commandArg)
	if err != nil || tryDelay < 1 || tryDelay > 1000000 {
		commandRequest.request.answer("Неправильный аргумент, отправьте `/delay N`, где N любое натуральное число меньше 1000000")
		return nil
	}
	commandRequest.request.customDelay.Delay = tryDelay
	if _, err := datastoreClient.Put(commandRequest.request.ctx, commandRequest.request.customDelayKey, &commandRequest.request.customDelay); err != nil {
		commandRequest.request.answerErrorWithLog("Не удалось сохранить, отправьте еще раз `/delay N`, где N любое натуральное число меньше 1000000", err)
		return nil
	}
	customDelayMap[commandRequest.request.updateMessage.Chat.ID] = commandRequest.request.customDelay.Delay
	commandRequest.request.answer("Я буду пропускать случайное число сообщений от 0 до " + commandArg)
	delete(delayMap, commandRequest.request.updateMessage.Chat.ID)
	return nil
}

func switchDatastoreBool(request *requestInfo, dsName string, value bool) {
	var localCache map[int64]bool
	var resultStruct *DatastoreBool
	var dsKey *datastore.Key
	switch dsName {
	case "Gentle":
		localCache = gentleMap
		resultStruct = &request.gentleStruct
		dsKey = request.gentleKey
	case "Stopped":
		localCache = stoppedMap
		resultStruct = &request.stoppedStruct
		dsKey = request.stoppedKey
	default:
		return
	}

	localCache[request.updateMessage.Chat.ID] = value
	resultStruct.Value = value
	if _, err := datastoreClient.Put(request.ctx, dsKey, resultStruct); err != nil {
		request.logWarn(err)
	}
}

type commandHardcore botCommand

func (commandRequest *commandHardcore) Handle() error {
	switchDatastoreBool(commandRequest.request, "Gentle", false)
	commandRequest.request.answer("Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle")
	return nil
}

type commandGentle botCommand

func (commandRequest *commandGentle) Handle() error {
	switchDatastoreBool(commandRequest.request, "Gentle", true)
	commandRequest.request.answer("Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore")
	return nil
}

type commandAmount botCommand

func (commandRequest *commandAmount) Handle() error {
	command := strings.Fields(commandRequest.request.updateMessage.Text)
	if len(command) < 2 {
		currentWordsAmount := 1
		if currentAmount, ok := wordsAmountMap[commandRequest.request.updateMessage.Chat.ID]; ok {
			currentWordsAmount = currentAmount
		}
		commandRequest.request.answer("Сейчас я хуифицирую случайное число слов от 1 до " + strconv.Itoa(currentWordsAmount))
		return nil
	}
	commandArg := command[len(command)-1]
	tryWordsAmount, err := strconv.Atoi(commandArg)
	if err != nil || tryWordsAmount < 1 || tryWordsAmount > 10 {
		commandRequest.request.answer("Неправильный аргумент, отправьте `/amount N`, где N любое натуральное число не больше 10")
		return nil
	}
	commandRequest.request.wordsAmountStruct.Value = tryWordsAmount
	if _, err := datastoreClient.Put(commandRequest.request.ctx, commandRequest.request.wordsAmountKey, &commandRequest.request.wordsAmountStruct); err != nil {
		commandRequest.request.answerErrorWithLog("Не удалось сохранить, отправьте еще раз `/amount N`, где N любое натуральное число не больше 10", err)
		return nil
	}
	wordsAmountMap[commandRequest.request.updateMessage.Chat.ID] = commandRequest.request.wordsAmountStruct.Value
	commandRequest.request.answer("Я буду хуифицировать случайное число слов от 1 до " + strconv.Itoa(commandRequest.request.wordsAmountStruct.Value))
	return nil
}

type commandNotFound botCommand

func (commandRequest *commandNotFound) Handle() error {
	return errors.New("Команда не найдена")
}
