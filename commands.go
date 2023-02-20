package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type commandStart botCommand

func (commandRequest *commandStart) Handle() error {
	message := "Я Хуебот\n" +
		"Я буду хуифицировать некоторые из ваших фраз\n" +
		"Сейчас режим вежливости *%s*\n" +
		"За подробностями в /help"
	settings.cache[commandRequest.request.cacheID].Enabled = true
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.logWarn(err)
		// Do not send error to command
		return nil
	}
	if settings.cache[commandRequest.request.cacheID].Gentle {
		message = fmt.Sprintf(message, "включен")
	} else {
		message = fmt.Sprintf(message, "отключен")
	}
	commandRequest.request.answer(message, MarkdownV2)
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
	commandRequest.request.answer("Выключаюсь", "")
	return nil
}

type commandHelp botCommand

func (commandRequest *commandHelp) Handle() error {
	commandRequest.request.answer("Статус бота: *"+commandRequest.request.getStatusString()+"*\n\n"+
		"Настройки:\n"+
		"  Вежливый режим:\n"+
		"    Включить:  /gentle\n"+
		"    Выключить: /hardcore\n"+
		"  Режим ответов:\n"+
		"    Включить: /reply\n"+
		"    Выключить: /noreply\n"+
		"  Частота ответов: `/delay N`, где _N_ от `1` до `500`\n"+
		"  Число хуифицируемых слов: `/amount N`, где _N_ от `1` до `10`\n"+
		"  Для остановки используйте /stop\n"+
		"  Для перезапуска используйте /start\n\n"+
		"По вопросам пишите на [GitHub](https://github.com/m-messiah/xye-bot/issues/new)\n"+
		"[Donate](https://www.paypal.com/donate?hosted_button_id=7KT6MFDSHPXL6)", MarkdownV2)
	return nil
}

type commandDelay botCommand

func (commandRequest *commandDelay) Handle() error {
	command := strings.Fields(commandRequest.request.updateMessage.Text)
	if len(command) < 2 {
		commandRequest.request.answer("Сейчас я пропускаю случайное число сообщений от `0` до `"+strconv.Itoa(settings.cache[commandRequest.request.cacheID].Delay)+"`", MarkdownV2)
		return nil
	}
	commandArg := command[len(command)-1]
	tryDelay, err := strconv.Atoi(commandArg)
	if err != nil || tryDelay < 1 || tryDelay > 500 {
		commandRequest.request.answer("Неправильный аргумент, отправьте `/delay N`, где _N_ любое натуральное число меньше `500`", MarkdownV2)
		return nil
	}
	settings.cache[commandRequest.request.cacheID].Delay = tryDelay
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.answerErrorWithLog("Не удалось сохранить, отправьте еще раз `/delay N`, где _N_ любое натуральное число меньше `500`", err, MarkdownV2)
		return nil
	}
	commandRequest.request.answer("Я буду пропускать случайное число сообщений от `0` до `"+strconv.Itoa(settings.cache[commandRequest.request.cacheID].Delay)+"`", MarkdownV2)
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
	commandRequest.request.answer("Вежливый режим *отключен*\nЧтобы включить его, используйте команду /gentle", MarkdownV2)
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
	commandRequest.request.answer("Вежливый режим *включен*\nЧтобы отключить его, используйте команду /hardcore", MarkdownV2)
	return nil
}

type commandReply botCommand

func (commandRequest *commandReply) Handle() error {
	settings.cache[commandRequest.request.cacheID].Reply = true
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.logWarn(err)
		// Do not send error to command
		return nil
	}
	commandRequest.request.answer("Режим ответов на сообщения *включен*\nЧтобы отключить его, используйте команду /noreply", MarkdownV2)
	return nil
}

type commandNoReply botCommand

func (commandRequest *commandNoReply) Handle() error {
	settings.cache[commandRequest.request.cacheID].Reply = false
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.logWarn(err)
		// Do not send error to command
		return nil
	}
	commandRequest.request.answer("Режим ответов на сообщения *отключен*\nЧтобы включить его, используйте команду /reply", MarkdownV2)
	return nil
}

type commandAmount botCommand

func (commandRequest *commandAmount) Handle() error {
	command := strings.Fields(commandRequest.request.updateMessage.Text)
	if len(command) < 2 {
		commandRequest.request.answer("Сейчас я хуифицирую случайное число слов от `1` до `"+strconv.Itoa(settings.cache[commandRequest.request.cacheID].WordsAmount)+"`", MarkdownV2)
		return nil
	}
	commandArg := command[len(command)-1]
	tryWordsAmount, err := strconv.Atoi(commandArg)
	if err != nil || tryWordsAmount < 1 || tryWordsAmount > 10 {
		commandRequest.request.answer("Неправильный аргумент, отправьте `/amount N`, где _N_ любое натуральное число не больше `10`", MarkdownV2)
		return nil
	}
	settings.cache[commandRequest.request.cacheID].WordsAmount = tryWordsAmount
	if err := settings.SaveCache(commandRequest.request.ctx, commandRequest.request.cacheID); err != nil {
		commandRequest.request.answerErrorWithLog("Не удалось сохранить, отправьте еще раз `/amount N`, где _N_ любое натуральное число не больше `10`", err, MarkdownV2)
		return nil
	}
	commandRequest.request.answer("Я буду хуифицировать случайное число слов от `1` до `"+strconv.Itoa(settings.cache[commandRequest.request.cacheID].WordsAmount)+"`", MarkdownV2)
	return nil
}

type commandNotFound botCommand

func (commandRequest *commandNotFound) Handle() error {
	return errors.New("команда не найдена")
}
