package main

import (
	"context"
	"encoding/json"
	"fmt"
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

// DatastoreDelay type for DataStore
type DatastoreDelay struct {
	Delay int
}

// DatastoreBool type for DataStore
type DatastoreBool struct {
	Value bool
}

type DatastoreGentle struct {
	Gentle bool
	Value  bool
}

// DatastoreInt type for DataStore
type DatastoreInt struct {
	Value int
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

func findIndex(keys []*datastore.Key, key int64) int {
	for i, k := range keys {
		if k.ID == key {
			return i
		}
	}
	return -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func migrate() {
	stoppedValues := make([]DatastoreBool, 200000)
	stoppedKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("Stopped"), &stoppedValues)
	if err != nil {
		log.Printf("unable to get Stopped keys: %s", err)
		return
	}
	log.Printf("got %d stopped keys", len(stoppedKeys))
	gentleValues := make([]DatastoreGentle, 200000)
	gentleKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("Gentle"), &gentleValues)
	if err != nil {
		log.Printf("unable to get Gentle keys: %s", err)
		return
	}
	log.Printf("got %d gentle keys", len(gentleKeys))
	delayValues := make([]DatastoreDelay, 200000)
	delayKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("DatastoreDelay"), &delayValues)
	if err != nil {
		log.Printf("unable to get Delay keys: %s", err)
		return
	}
	log.Printf("got %d delay keys", len(delayKeys))
	wordsValues := make([]DatastoreInt, 200000)
	wordsKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("WordsAmount"), &wordsValues)
	if err != nil {
		log.Printf("unable to get WordsAmount keys: %s", err)
		return
	}
	log.Printf("got %d words keys", len(wordsKeys))

	// log.Printf("Cleanup chatSettings")
	// chatSettingsKeys, err := settings.client.GetAll(context.Background(), datastore.NewQuery("ChatSettings").KeysOnly(), nil)
	// for _, key := range chatSettingsKeys {
	// 	log.Printf("delete %s", key)
	// 	settings.client.Delete(context.Background(), key)
	// }
	log.Printf("chatSettings cleaned up. Starting migration")
	for keyIndex, stoppedKey := range stoppedKeys {
		chatSettings := settings.DefaultChatSettings()
		chatSettings.Enabled = !stoppedValues[keyIndex].Value
		if i := findIndex(gentleKeys, stoppedKey.ID); i > -1 {
			chatSettings.Gentle = gentleValues[i].Gentle || gentleValues[i].Value
		}
		if i := findIndex(delayKeys, stoppedKey.ID); i > -1 {
			chatSettings.Delay = min(delayValues[i].Delay, 0)
		}
		if i := findIndex(wordsKeys, stoppedKey.ID); i > -1 {
			chatSettings.WordsAmount = min(wordsValues[i].Value, 1)
		}
		settings.cache[fmt.Sprintf("%d", stoppedKey.ID)] = &chatSettings
		if err := settings.SaveCache(context.Background(), fmt.Sprintf("%d", stoppedKey.ID)); err != nil {
			log.Printf("could not save %s (%v): %s", fmt.Sprintf("%d", stoppedKey.ID), chatSettings, err)
		}
		log.Printf("saved successfully %d", stoppedKey.ID)
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
