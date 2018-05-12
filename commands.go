package xyebot

import (
	"errors"
	"fmt"
	"google.golang.org/appengine/datastore"
	"strconv"
	"strings"
)

type CommandStart Command

func (self *CommandStart) Handle() error {
	message := "Привет! Я бот-хуебот.\n" +
		"Я буду хуифицировать некоторые из ваших фраз.\n" +
		"Сейчас режим вежливости %s\n" +
		"За подробностями в /help"
	switchDatastoreBool(self.request, "Stopped", false)
	if Gentle[self.request.updateMessage.Chat.ID] {
		message = fmt.Sprintf(message, "включен")
	} else {
		message = fmt.Sprintf(message, "отключен")
	}
	self.request.Answer(message)
	return nil
}

type CommandStop Command

func (self *CommandStop) Handle() error {
	switchDatastoreBool(self.request, "Stopped", true)
	self.request.Answer("Выключаюсь")
	return nil
}

type CommandHelp Command

func (self *CommandHelp) Handle() error {
	self.request.Answer(
		"Вежливый режим:\n" +
			"  Для включения используйте команду /gentle\n" +
			"  Для отключения - /hardcore\n" +
			"Частота ответов: /delay N, где N - любое любое натуральное число\n" +
			"Число хуифицируемых слов: /amount N, где N - от 1 до 10\n" +
			"Для остановки используйте /stop")
	return nil
}

type CommandDelay Command

func (self *CommandDelay) Handle() error {
	command := strings.Fields(self.request.updateMessage.Text)
	if len(command) < 2 {
		currentDelayMessage := "Сейчас я пропускаю случайное число сообщений от 0 до "
		if currentDelay, ok := CustomDelay[self.request.updateMessage.Chat.ID]; ok {
			currentDelayMessage += strconv.Itoa(currentDelay)
		} else {
			currentDelayMessage += "4"
		}
		self.request.Answer(currentDelayMessage)
		return nil
	}
	commandArg := command[len(command)-1]
	tryDelay, err := strconv.Atoi(commandArg)
	if err != nil || tryDelay < 1 || tryDelay > 1000000 {
		self.request.Answer("Неправильный аргумент, отправьте `/delay N`, где N любое натуральное число меньше 1000000")
		return nil
	}
	self.request.customDelay.Delay = tryDelay
	if _, err := datastore.Put(self.request.ctx, self.request.customDelayKey, &self.request.customDelay); err != nil {
		self.request.AnswerErrorWithLog("Не удалось сохранить, отправьте еще раз `/delay N`, где N любое натуральное число меньше 1000000", err)
		return nil
	}
	CustomDelay[self.request.updateMessage.Chat.ID] = self.request.customDelay.Delay
	self.request.Answer("Я буду пропускать случайное число сообщений от 0 до " + commandArg)
	delete(Delay, self.request.updateMessage.Chat.ID)
	return nil
}

func switchDatastoreBool(request *Request, dsName string, value bool) {
	var localCache map[int64]bool
	var resultStruct *DatastoreBool
	var dsKey *datastore.Key
	switch dsName {
	case "Gentle":
		localCache = Gentle
		resultStruct = &request.gentleStruct
		dsKey = request.gentleKey
	case "Stopped":
		localCache = Stopped
		resultStruct = &request.stoppedStruct
		dsKey = request.stoppedKey
	default:
		return
	}

	localCache[request.updateMessage.Chat.ID] = value
	resultStruct.Value = value
	if _, err := datastore.Put(request.ctx, dsKey, &resultStruct); err != nil {
		request.LogWarn(err)
	}
}

type CommandHardcore Command

func (self *CommandHardcore) Handle() error {
	switchDatastoreBool(self.request, "Gentle", false)
	self.request.Answer("Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle")
	return nil
}

type CommandGentle Command

func (self *CommandGentle) Handle() error {
	switchDatastoreBool(self.request, "Gentle", true)
	self.request.Answer("Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore")
	return nil
}

type CommandAmount Command

func (self *CommandAmount) Handle() error {
	command := strings.Fields(self.request.updateMessage.Text)
	if len(command) < 2 {
		currentWordsAmount := 1
		if currentAmount, ok := WordsAmount[self.request.updateMessage.Chat.ID]; ok {
			currentWordsAmount = currentAmount
		}
		self.request.Answer("Сейчас я хуифицирую случайное число слов от 1 до " + strconv.Itoa(currentWordsAmount))
		return nil
	}
	commandArg := command[len(command)-1]
	tryWordsAmount, err := strconv.Atoi(commandArg)
	if err != nil || tryWordsAmount < 1 || tryWordsAmount > 10 {
		self.request.Answer("Неправильный аргумент, отправьте `/amount N`, где N любое натуральное число не больше 10")
		return nil
	}
	self.request.wordsAmountStruct.Value = tryWordsAmount
	if _, err := datastore.Put(self.request.ctx, self.request.wordsAmountKey, &self.request.wordsAmountStruct); err != nil {
		self.request.AnswerErrorWithLog("Не удалось сохранить, отправьте еще раз `/amount N`, где N любое натуральное число не больше 10", err)
		return nil
	}
	WordsAmount[self.request.updateMessage.Chat.ID] = self.request.wordsAmountStruct.Value
	self.request.Answer("Я буду хуифицировать случайное число слов от 1 до " + strconv.Itoa(self.request.wordsAmountStruct.Value))
	return nil
}

type CommandNotFound Command

func (self *CommandNotFound) Handle() error {
	return errors.New("Команда не найдена")
}
