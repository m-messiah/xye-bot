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

func deleteKeys(name string, filtered bool) {
	var keysToDelete []*datastore.Key
	var err error
	if filtered {
		keysToDelete, err = settings.client.GetAll(context.Background(), datastore.NewQuery(name).Filter("Enabled = ", false).KeysOnly(), nil)
	} else {
		keysToDelete, err = settings.client.GetAll(context.Background(), datastore.NewQuery(name).KeysOnly(), nil)
	}
	if err != nil {
		log.Printf("unable to get %s keys: %s", name, err)
		return
	}
	log.Printf("got %s keys: %d", name, len(keysToDelete))
	batch := 400
	for i := 0; i < len(keysToDelete); i += batch {
		j := i + batch
		if j >= len(keysToDelete) {
			j = len(keysToDelete) - 1
		}
		if err := settings.client.DeleteMulti(context.Background(), keysToDelete[i:j]); err != nil {
			log.Printf("unable to delete %s keys: %s", name, err)
			return
		}
	}
}

func migrate() {
	deleteKeys("Stopped", false)
	deleteKeys("Gentle", false)
	deleteKeys("DatastoreDelay", false)
	deleteKeys("WordsAmount", false)
	deleteKeys("ChatSettings", true)
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
