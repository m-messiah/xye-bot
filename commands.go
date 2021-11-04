package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type commandStatus botCommand

func (commandRequest *commandStatus) Handle() error {
	commandRequest.request.answer(fmt.Sprintf("Processed items: %d", len(settings.cache)))
	return nil
}

type commandStop botCommand

func (commandRequest *commandStop) Handle() error {
	settings.cache[commandRequest.request.cacheID].Enabled = false
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.logWarn(err)
		// Do not send error to command
		return nil
	}
	commandRequest.request.answer("Выключаюсь")
	return nil
}

type commandHelp botCommand

func (commandRequest *commandHelp) Handle() error {
	commandRequest.request.answer(
		"Статус бота: " + commandRequest.request.getStatusString() + "\n\n" +
			"* Вежливый режим:\n" +
			"  * Для включения используйте команду /gentle\n" +
			"  * Для отключения - /hardcore\n" +
			"* Частота ответов: /delay N, где N - любое любое натуральное число\n" +
			"* Число хуифицируемых слов: /amount N, где N - от 1 до 10\n" +
			"* Для остановки используйте /stop\n" +
			"* Для перезапуска используйте /start\n\n" +
			"По вопросам: https://github.com/m-messiah/xye-bot/issues")
	return nil
}

type commandDelay botCommand

func (commandRequest *commandDelay) Handle() error {
	command := strings.Fields(commandRequest.request.updateMessage.Text)
	if len(command) < 2 {
		commandRequest.request.answer("Сейчас я пропускаю случайное число сообщений от 0 до " + strconv.Itoa(settings.cache[commandRequest.request.cacheID].Delay))
		return nil
	}
	commandArg := command[len(command)-1]
	tryDelay, err := strconv.Atoi(commandArg)
	if err != nil || tryDelay < 1 || tryDelay > 1000000 {
		commandRequest.request.answer("Неправильный аргумент, отправьте `/delay N`, где N любое натуральное число меньше 1000000")
		return nil
	}
	settings.cache[commandRequest.request.cacheID].Delay = tryDelay
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.answerErrorWithLog("Не удалось сохранить, отправьте еще раз `/delay N`, где N любое натуральное число меньше 1000000", err)
		return nil
	}
	commandRequest.request.answer("Я буду пропускать случайное число сообщений от 0 до " + strconv.Itoa(settings.cache[commandRequest.request.cacheID].Delay))
	delete(delayMap, commandRequest.request.updateMessage.Chat.ID)
	return nil
}

type commandHardcore botCommand

func (commandRequest *commandHardcore) Handle() error {
	settings.cache[commandRequest.request.cacheID].Gentle = false
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.logWarn(err)
		// Do not send error to command
		return nil
	}
	commandRequest.request.answer("Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle")
	return nil
}

type commandGentle botCommand

func (commandRequest *commandGentle) Handle() error {
	settings.cache[commandRequest.request.cacheID].Gentle = true
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.logWarn(err)
		// Do not send error to command
		return nil
	}
	commandRequest.request.answer("Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore")
	return nil
}

type commandAmount botCommand

func (commandRequest *commandAmount) Handle() error {
	command := strings.Fields(commandRequest.request.updateMessage.Text)
	if len(command) < 2 {
		commandRequest.request.answer("Сейчас я хуифицирую случайное число слов от 1 до " + strconv.Itoa(settings.cache[commandRequest.request.cacheID].WordsAmount))
		return nil
	}
	commandArg := command[len(command)-1]
	tryWordsAmount, err := strconv.Atoi(commandArg)
	if err != nil || tryWordsAmount < 1 || tryWordsAmount > 10 {
		commandRequest.request.answer("Неправильный аргумент, отправьте `/amount N`, где N любое натуральное число не больше 10")
		return nil
	}
	settings.cache[commandRequest.request.cacheID].WordsAmount = tryWordsAmount
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.answerErrorWithLog("Не удалось сохранить, отправьте еще раз `/amount N`, где N любое натуральное число не больше 10", err)
		return nil
	}
	commandRequest.request.answer("Я буду хуифицировать случайное число слов от 1 до " + strconv.Itoa(settings.cache[commandRequest.request.cacheID].WordsAmount))
	return nil
}

type commandNotFound botCommand

func (commandRequest *commandNotFound) Handle() error {
	return errors.New("команда не найдена")
}
