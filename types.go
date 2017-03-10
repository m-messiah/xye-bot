package xyebot

type Response struct {
	Chatid int64  `json:"chat_id"`
	Text   string `json:"text"`
	Method string `json:"method"`
}

type Chat struct {
	ID int64 `json: "chat_id"`
}

type Message struct {
	Chat *Chat  `json:"chat"`
	Text string `json:"text"`
}

type Update struct {
	Message *Message `json:"message"`
}

type DatastoreDelay struct {
	Delay int
}

type DatastoreGentle struct {
	Gentle bool
}
