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

var (
	delayMap map[int64]int
	settings Settings
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
}

func migrate() {
	stoppedKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("Stopped").KeysOnly(), nil)
	if err != nil {
		log.Printf("unable to get Stopped keys: %s", err)
		return
	}
	if err := settings.client.DeleteMulti(context.Background(), stoppedKeys); err != nil {
		log.Printf("unable to delete Stopped keys: %s", err)
		return
	}

	gentleKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("Gentle").KeysOnly(), nil)
	if err != nil {
		log.Printf("unable to get Gentle keys: %s", err)
		return
	}
	if err := settings.client.DeleteMulti(context.Background(), gentleKeys); err != nil {
		log.Printf("unable to delete Gentle keys: %s", err)
		return
	}

	delayKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("DatastoreDelay").KeysOnly(), nil)
	if err != nil {
		log.Printf("unable to get Delay keys: %s", err)
		return
	}
	if err := settings.client.DeleteMulti(context.Background(), delayKeys); err != nil {
		log.Printf("unable to delete Delay keys: %s", err)
		return
	}

	wordsKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("WordsAmount").KeysOnly(), nil)
	if err != nil {
		log.Printf("unable to get WordsAmount keys: %s", err)
		return
	}
	if err := settings.client.DeleteMulti(context.Background(), wordsKeys); err != nil {
		log.Printf("unable to delete WordsAmount keys: %s", err)
		return
	}

	settingsKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("ChatSettings").Filter("Enabled = ", false).KeysOnly(), nil)
	if err != nil {
		log.Printf("unable to get ChatSettings keys: %s", err)
		return
	}
	if err := settings.client.DeleteMulti(context.Background(), settingsKeys); err != nil {
		log.Printf("unable to delete ChatSettings disabled keys: %s", err)
		return
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	settings = NewSettings()
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	go migrate()

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
