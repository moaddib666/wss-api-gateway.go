package protocol

import "encoding/json"

type Message struct {
	// TODO introduce message id
	Metadata *Metadata `json:"metadata"`
	Payload  []byte    `json:"payload"`
}

func (m *Message) Dump() ([]byte, error) {
	return json.Marshal(m)
}

func NewMessage(sender, recipient string, payload []byte) *Message {
	obj := &Message{
		Metadata: &Metadata{
			Sender:    sender,
			Recipient: recipient,
		},
		Payload: payload,
	}
	return obj
}

type Metadata struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}
