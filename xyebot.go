package xyebot

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

const DEFAULT_DELAY = 4

var Delay map[int64]int
var Gentle map[int64]bool
var WordsAmount map[int64]int
var Stopped map[int64]bool
var CustomDelay map[int64]int

func SendMessage(w http.ResponseWriter, chatID int64, text string, replyToID *int64) {
	var msg Response
	if replyToID == nil {
		msg = Response{Chatid: chatID, Text: text, Method: "sendMessage"}
	} else {
		msg = Response{Chatid: chatID, Text: text, ReplyToID: replyToID, Method: "sendMessage"}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	request, err := NewRequest(w, r)
	if err != nil {
		return
	}
	if err = request.ParseCommand(w); err == nil {
		return
	}
	if request.IsStopped() {
		return
	}
	request.HandleDelay()
	replyID := request.GetReplyIDIfNeeded()
	if request.IsAnswerNeeded(replyID) {
		if replyID == nil {
			request.CleanDelay()
		}
		// log.Infof(ctx, "[%v] %s", updateMessage.Chat.ID, updateMessage.Text)
		output := request.Huify()
		if output != "" {
			SendMessage(request.writer, request.updateMessage.Chat.ID, output, replyID)
			return
		}
	}
}

func init() {
	Delay = make(map[int64]int)
	Gentle = make(map[int64]bool)
	WordsAmount = make(map[int64]int)
	Stopped = make(map[int64]bool)
	CustomDelay = make(map[int64]int)
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", Handler)
}
