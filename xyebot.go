package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const (
	// BotName is the telegram name of the bot. Needs for commands identification
	BotName = "xye_bot"
	// DelayLimit is the upper limit of randomly skipped messages.
	// For example, 4 means bot would skip random amount of messages between 1 and 4.
	DelayLimit = 4
	// WordsAmount is the amount of words from the message to apply modification
	WordsAmount = 1
)

var delayMap map[int64]int

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
	if err = request.handleCommand(); err == nil {
		return
	}
	request.handleDelay()
	if !request.isAnswerNeeded() {
		return
	}
	output := request.Modify()
	if output == "" {
		return
	}
	request.cleanDelay()
	sendMessage(request.writer, request.updateMessage.Chat.ID, output, &request.updateMessage.ID, "")
}

func main() {
	delayMap = make(map[int64]int)
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
