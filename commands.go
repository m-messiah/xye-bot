package xyebot

import (
	"errors"
	"fmt"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"strconv"
	"strings"
)

type CommandStart Command

func (self *CommandStart) Handle() error {
	message := "Привет! Я бот-хуебот.\n" +
		"Я буду хуифицировать некоторые из Ваших фраз.\n" +
		"Сейчас режим вежливости %s\n" +
		"За подробностями в /help"
	Stopped[self.request.updateMessage.Chat.ID] = false
	self.request.stoppedStruct.Value = false
	if _, err := datastore.Put(self.request.ctx, self.request.stoppedKey, &self.request.stoppedStruct); err != nil {
		log.Warningf(self.request.ctx, "[%v] %s", self.request.updateMessage.Chat.ID, err.Error())
	}
	if Gentle[self.request.updateMessage.Chat.ID] {
		message = fmt.Sprintf(message, "включен")
	} else {
		message = fmt.Sprintf(message, "отключен")
	}
	SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, message, nil)
	return nil
}

type CommandStop Command

func (self *CommandStop) Handle() error {
	Stopped[self.request.updateMessage.Chat.ID] = true
	self.request.stoppedStruct.Value = true
	if _, err := datastore.Put(self.request.ctx, self.request.stoppedKey, &self.request.stoppedStruct); err != nil {
		log.Warningf(self.request.ctx, "[%v] %s", self.request.updateMessage.Chat.ID, err.Error())
	}
	SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Выключаюсь", nil)
	return nil
}

type CommandHelp Command

func (self *CommandHelp) Handle() error {
	SendMessage(self.request.writer, self.request.updateMessage.Chat.ID,
		"Вежливый режим:\n"+
			"  Для включения используйте команду /gentle\n"+
			"  Для отключения - /hardcore\n"+
			"Частота ответов: /delay N, где N - любое любое натуральное число\n"+
			"Число хуифицируемых слов: /amount N, где N - от 1 до 10\n"+
			"Для остановки используйте /stop", nil)
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
		SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, currentDelayMessage, nil)
		return nil
	}
	commandArg := command[len(command)-1]
	tryDelay, err := strconv.Atoi(commandArg)
	if err != nil || tryDelay < 1 || tryDelay > 1000000 {
		SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Неправильный аргумент, отправьте `/delay N`, где N любое натуральное число меньше 1000000", nil)
		return nil
	}
	self.request.customDelay.Delay = tryDelay
	if _, err := datastore.Put(self.request.ctx, self.request.customDelayKey, &self.request.customDelay); err != nil {
		log.Warningf(self.request.ctx, "[%v] %s", self.request.updateMessage.Chat.ID, err.Error())
		SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Не удалось сохранить, отправьте еще раз `/delay N`, где N любое натуральное число меньше 1000000", nil)
		return nil
	}
	CustomDelay[self.request.updateMessage.Chat.ID] = self.request.customDelay.Delay
	SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Я буду пропускать случайное число сообщений от 0 до "+commandArg, nil)
	delete(Delay, self.request.updateMessage.Chat.ID)
	return nil
}

type CommandHardcore Command

func (self *CommandHardcore) Handle() error {
	Gentle[self.request.updateMessage.Chat.ID] = false
	self.request.gentleStruct.Value = false
	if _, err := datastore.Put(self.request.ctx, self.request.gentleKey, &self.request.gentleStruct); err != nil {
		log.Warningf(self.request.ctx, "[%v] %s", self.request.updateMessage.Chat.ID, err.Error())
	}
	SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle", nil)
	return nil
}

type CommandGentle Command

func (self *CommandGentle) Handle() error {
	Gentle[self.request.updateMessage.Chat.ID] = true
	self.request.gentleStruct.Value = true
	if _, err := datastore.Put(self.request.ctx, self.request.gentleKey, &self.request.gentleStruct); err != nil {
		log.Warningf(self.request.ctx, "[%v] %s", self.request.updateMessage.Chat.ID, err.Error())
	}
	SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore", nil)
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
		SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Сейчас я хуифицирую случайное число слов от 1 до "+strconv.Itoa(currentWordsAmount), nil)
		return nil
	}
	commandArg := command[len(command)-1]
	tryWordsAmount, err := strconv.Atoi(commandArg)
	if err != nil || tryWordsAmount < 1 || tryWordsAmount > 10 {
		SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Неправильный аргумент, отправьте `/amount N`, где N любое натуральное число не больше 10", nil)
		return nil
	}
	self.request.wordsAmountStruct.Value = tryWordsAmount
	if _, err := datastore.Put(self.request.ctx, self.request.wordsAmountKey, &self.request.wordsAmountStruct); err != nil {
		log.Warningf(self.request.ctx, "[%v] %s", self.request.updateMessage.Chat.ID, err.Error())
		SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Не удалось сохранить, отправьте еще раз `/amount N`, где N любое натуральное число не больше 10", nil)
		return nil
	}
	WordsAmount[self.request.updateMessage.Chat.ID] = self.request.wordsAmountStruct.Value
	SendMessage(self.request.writer, self.request.updateMessage.Chat.ID, "Я буду хуифицировать случайное число слов от 1 до "+strconv.Itoa(self.request.wordsAmountStruct.Value), nil)
	return nil
}

type CommandNotFound Command

func (self *CommandNotFound) Handle() error {
	return errors.New("Команда не найдена")
}
