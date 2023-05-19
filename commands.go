package main

import (
	"errors"
)

type commandStart botCommand

func (commandRequest *commandStart) Handle() error {
	commandRequest.request.answer("Я вежливый Хуебот\nЯ буду хуифицировать некоторые из ваших фраз", "")
	return nil
}

type commandStop botCommand

func (commandRequest *commandStop) Handle() error {
	commandRequest.request.answer("Я не умею обрабатывать эту команду. Удалите меня из чата, если надоел", "")
	return nil
}

type commandHelp botCommand

func (commandRequest *commandHelp) Handle() error {
	commandRequest.request.answer("Бот использует вежливый режим и отвечает на случайные сообщения\n\n"+
		"Если хочется запустить своего бота, код на [GitHub](https://github.com/m-messiah/xye-bot)\n"+
		"[Donate](https://www.paypal.com/donate?hosted_button_id=7KT6MFDSHPXL6)", "MarkdownV2")
	return nil
}

type commandNotFound botCommand

func (commandRequest *commandNotFound) Handle() error {
	return errors.New("команда не найдена")
}
