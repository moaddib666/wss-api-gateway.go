package protocol

import (
	"WSSFacade/constants"
	"encoding/json"
	"time"
)

type Event struct {
	Metadata struct {
		EntityId      string `json:"entity_id"`
		EntityName    string `json:"entity_name"`
		PublisherName string `json:"publisher_name"`
		EventName     string `json:"event_name"`
		Created       string `json:"created"`
	} `json:"metadata"`
	Payload struct{} `json:"payload"`
}

func (e Event) Dump() []byte {
	data, _ := json.Marshal(e)
	return data
}

func NewClientEvent(clientId string, eventName string) *Event {
	metadata := struct {
		EntityId      string `json:"entity_id"`
		EntityName    string `json:"entity_name"`
		PublisherName string `json:"publisher_name"`
		EventName     string `json:"event_name"`
		Created       string `json:"created"`
	}{
		clientId,
		"client",
		constants.AppInternalName,
		eventName,
		time.Now().Format(time.RFC3339),
	}
	payload := struct{}{}
	return &Event{Metadata: metadata, Payload: payload}
}

func ClientConnectedEvent(clientId string) *Event {
	return NewClientEvent(clientId, "connected")
}

func ClientDisconnectedEvent(clientId string) *Event {
	return NewClientEvent(clientId, "disconnected")
}
