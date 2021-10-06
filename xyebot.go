package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

const defaultDelay = 4

var (
	delayMap        map[int64]int
	gentleMap       map[int64]bool
	wordsAmountMap  map[int64]int
	stoppedMap      map[int64]bool
	customDelayMap  map[int64]int
	datastoreClient *datastore.Client
)

func sendMessage(w http.ResponseWriter, chatID int64, text string, replyToID *int64) {
	var msg Response
	if replyToID == nil {
		msg = Response{ChatID: chatID, Text: text, Method: "sendMessage"}
	} else {
		msg = Response{ChatID: chatID, Text: text, ReplyToID: replyToID, Method: "sendMessage"}
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
			sendMessage(request.writer, request.updateMessage.Chat.ID, output, replyID)
			return
		}
	}
}

func main() {
	delayMap = make(map[int64]int)
	gentleMap = make(map[int64]bool)
	wordsAmountMap = make(map[int64]int)
	stoppedMap = make(map[int64]bool)
	customDelayMap = make(map[int64]int)
	rand.Seed(time.Now().UTC().UnixNano())
	var err error
	datastoreClient, err = datastore.NewClient(context.Background(), "xye-bot")
	if err != nil {
		log.Fatal(err)
	}

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
