package xyebot

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var Delay map[int64]int
var Gentle map[int64]bool
var WordsAmount map[int64]int
var Stopped map[int64]bool
var CustomDelay map[int64]int

const DEFAULT_DELAY = 4

func sendMessage(w http.ResponseWriter, chatID int64, text string, replyToID *int64) {
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

func isCommand(text, command string) bool {
	if strings.Index(text, command) == 0 {
		if strings.Contains(text, "@xye_bot") || !strings.Contains(text, "@") {
			return true
		}
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	request, err := newRequest(r)
	if err != nil {
		return
	}
	if err = request.parseCommand(w); err == nil {
		return
	}
	if request.isStopped() {
		return
	}
	request.handleDelay()
	replyID := request.getReplyIDIfNeeded()
	if request.isAnswerNeeded(replyID) {
		if replyID == nil {
			request.cleanDelay()
		}
		// log.Infof(ctx, "[%v] %s", updateMessage.Chat.ID, updateMessage.Text)
		output := request.huify()
		if output != "" {
			sendMessage(w, request.updateMessage.Chat.ID, output, replyID)
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

	http.HandleFunc("/", handler)
}
