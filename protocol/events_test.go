package protocol

import "testing"

func TestClientConnectedEvent(t *testing.T) {
	event := ClientConnectedEvent("test")
	if event.Meta.ObjectType != "client" {
		t.Errorf("Error creating client connected event: %s", event.Meta.ObjectType)
	}
	if event.Meta.Event != "connected" {
		t.Errorf("Error creating client connected event: %s", event.Meta.Event)
	}
}

func TestClientDisconnectedEvent(t *testing.T) {
	event := ClientDisconnectedEvent("test")
	if event.Meta.ObjectType != "client" {
		t.Errorf("Error creating client disconnected event: %s", event.Meta.ObjectType)
	}
	if event.Meta.Event != "disconnected" {
		t.Errorf("Error creating client disconnected event: %s", event.Meta.Event)
	}
}
