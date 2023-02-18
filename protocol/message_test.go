package protocol

import "testing"

func TestNewMessage(t *testing.T) {
	message := NewMessage("sender", "recipient", []byte("payload"))
	if message.Metadata.Sender != "sender" {
		t.Errorf("Error creating message: %s", message.Metadata.Sender)
	}
	if message.Metadata.Recipient != "recipient" {
		t.Errorf("Error creating message: %s", message.Metadata.Recipient)
	}
	if string(message.Payload) != "payload" {
		t.Errorf("Error creating message: %s", string(message.Payload))
	}
}

func TestMessage_Dump(t *testing.T) {
	message := NewMessage("sender", "recipient", []byte("payload"))
	dump, err := message.Dump()
	if err != nil {
		t.Errorf("Error dumping message: %v", err)
	}
	if string(dump) != `{"metadata":{"sender":"sender","recipient":"recipient"},"payload":"cGF5bG9hZA=="}` {
		t.Errorf("Error dumping message: %s", string(dump))
	}
}
