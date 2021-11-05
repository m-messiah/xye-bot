package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const MarkdownV2 = "MarkdownV2"

var (
	delayMap map[int64]int
	settings Settings
)

func sendMessage(w http.ResponseWriter, chatID int64, text string, replyToID *int64, parseMode string) {
	var msg Response
	if replyToID == nil {
		msg = Response{ChatID: chatID, Text: text, Method: "sendMessage", ParseMode: parseMode, DisableWebPagePreview: true}
	} else {
		msg = Response{ChatID: chatID, Text: text, ReplyToID: replyToID, Method: "sendMessage", ParseMode: parseMode, DisableWebPagePreview: true}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(msg)
}

func handler(w http.ResponseWriter, r *http.Request) {
	request, err := newRequest(w, r)
	if err != nil {
		return
	}
	if err = request.parseCommand(); err == nil {
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
		output := request.huify()
		if output != "" {
			sendMessage(request.writer, request.updateMessage.Chat.ID, output, replyID, "")
			return
		}
	}
}

func main() {
	delayMap = make(map[int64]int)
	rand.Seed(time.Now().UTC().UnixNano())
	settings = NewSettings()
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
