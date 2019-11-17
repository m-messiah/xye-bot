package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"net/http"
)

// Response to Telegram
type Response struct {
	Chatid    int64  `json:"chat_id"`
	Text      string `json:"text"`
	Method    string `json:"method"`
	ReplyToID *int64 `json:"reply_to_message_id"`
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

// DatastoreDelay type for DataStore
type DatastoreDelay struct {
	Delay int
}

// DatastoreBool type for DataStore
type DatastoreBool struct {
	Value bool
}

// DatastoreInt type for DataStore
type DatastoreInt struct {
	Value int
}

type requestInfo struct {
	customDelay       DatastoreDelay
	gentleStruct      DatastoreBool
	wordsAmountStruct DatastoreInt
	stoppedStruct     DatastoreBool
	updateMessage     *Message
	ctx               context.Context
	customDelayKey    *datastore.Key
	gentleKey         *datastore.Key
	stoppedKey        *datastore.Key
	wordsAmountKey    *datastore.Key
	writer            http.ResponseWriter
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
