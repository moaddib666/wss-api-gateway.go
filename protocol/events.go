package protocol

import (
	"WSSFacade/constants"
	"encoding/json"
	"time"
)

type Event struct {
	Meta struct {
		ObjectId   string `json:"object_id"`
		ObjectType string `json:"object_type"`
		Publisher  string `json:"publisher"`
		Event      string `json:"event"`
		Created    string `json:"created"`
	} `json:"meta"`
	Data struct{} `json:"data"`
}

func (e Event) Dump() []byte {
	data, _ := json.Marshal(e)
	return data
}

func NewClientEvent(clientId string, eventName string) *Event {
	meta := struct {
		ObjectId   string `json:"object_id"`
		ObjectType string `json:"object_type"`
		Publisher  string `json:"publisher"`
		Event      string `json:"event"`
		Created    string `json:"created"`
	}{
		clientId,
		"client",
		constants.AppInternalName,
		eventName,
		time.Now().Format(time.RFC3339),
	}
	data := struct{}{}
	return &Event{Meta: meta, Data: data}
}

func ClientConnectedEvent(clientId string) *Event {
	return NewClientEvent(clientId, "connected")
}

func ClientDisconnectedEvent(clientId string) *Event {
	return NewClientEvent(clientId, "disconnected")
}
