package main

import (
	"context"
	"net/http"
)

// Response to Telegram
type Response struct {
	ChatID                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	Method                string `json:"method"`
	ReplyToID             *int64 `json:"reply_to_message_id"`
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

// Chat Telegram structure
type Chat struct {
	ID       int64   `json:"id"`
	Username *string `json:"username"`
}

// Message Telegram structure
type Message struct {
	ID      int64    `json:"message_id"`
	Chat    *Chat    `json:"chat"`
	From    *Chat    `json:"from"`
	Text    string   `json:"text"`
	ReplyTo *Message `json:"reply_to_message"`
}

// Update - outer Telegram structure
type Update struct {
	Message       *Message `json:"message"`
	EditedMessage *Message `json:"edited_message"`
}

type requestInfo struct {
	updateMessage *Message
	ctx           context.Context
	cacheID       string
	writer        http.ResponseWriter
}

type botCommand struct {
	request *requestInfo
}

type commandInterface interface {
	Handle() error
}

func handleCommand(command commandInterface) error {
	return command.Handle()
}
