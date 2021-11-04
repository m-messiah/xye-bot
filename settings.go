package main

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
)

type ChatSettings struct {
	Delay       int
	Enabled     bool
	Gentle      bool
	WordsAmount int
}

type Settings struct {
	client *datastore.Client
	cache  map[string]*ChatSettings
}

func initClient() *datastore.Client {
	datastoreClient, err := datastore.NewClient(context.Background(), "xye-bot")
	if err != nil {
		log.Fatal(err)
	}
	return datastoreClient
}

func NewSettings() Settings {
	return Settings{
		client: initClient(),
		cache:  make(map[string]*ChatSettings),
	}
}

func (s Settings) Put(ctx context.Context, key *datastore.Key, src interface{}) (err error) {
	if _, err = s.client.Put(ctx, key, src); err != nil {
		s.client = initClient()
		_, err = s.client.Put(ctx, key, src)
	}
	return err
}

func (s Settings) Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error) {
	if err = s.client.Get(ctx, key, dst); err != nil {
		s.client = initClient()
		err = s.client.Get(ctx, key, dst)
	}
	return err
}

func (s Settings) datastoreKey(key string) *datastore.Key {
	return datastore.NameKey("ChatSettings", key, nil)
}

func (s Settings) DefaultChatSettings() ChatSettings {
	return ChatSettings{
		Delay:       4,
		Enabled:     true,
		Gentle:      true,
		WordsAmount: 1,
	}
}

func (s Settings) EnsureCache(ctx context.Context, key string) {
	if _, ok := s.cache[key]; !ok {
		datastoreKey := s.datastoreKey(key)
		var resultStruct ChatSettings
		if err := settings.Get(ctx, datastoreKey, &resultStruct); err != nil {
			resultStruct = s.DefaultChatSettings()
			if err := settings.Put(ctx, datastoreKey, resultStruct); err != nil {
				log.Printf("[%v] %s %+v - %s", key, datastoreKey, resultStruct, err.Error())
			}
		}
		s.cache[key] = &resultStruct
	}
}

func (s Settings) SaveCache(ctx context.Context, key string) error {
	return s.Put(ctx, s.datastoreKey(key), s.cache[key])
}

func (s Settings) SetDelay(ctx context.Context, key string, delay int) error {
	chatSettings := s.cache[key]
	chatSettings.Delay = delay
	s.cache[key] = chatSettings
	return s.SaveCache(ctx, key)
}
