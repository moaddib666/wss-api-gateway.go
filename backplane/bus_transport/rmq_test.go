package bus_transport

import (
	"MargayGateway/protocol"
	"testing"
)

// Note that in order for these tests to run successfully,
//	you will need a running instance of RabbitMQ on localhost:5672,
//	and the environment variable MARGAY_TRANSPORT_DSN should not be set.
//	If you have a different configuration, you may need to modify the tests accordingly.

func TestRMQTransport_SendMessage(t *testing.T) {
	transport := NewRMQTransport()
	err := transport.Init()
	if err != nil {
		t.Fatalf("error initializing RMQTransport: %v", err)
	}
	defer func() {
		err = transport.Destruct()
		if err != nil {
			t.Fatalf("error destructing RMQTransport: %v", err)
		}
	}()

	msg := &protocol.Message{
		Payload: []byte("test message"),
		Metadata: &protocol.Metadata{
			Sender:    "test-sender",
			Recipient: "test-recipient",
		},
	}

	err = transport.SendMessage(msg)
	if err != nil {
		t.Fatalf("error sending message: %v", err)
	}
}

func TestRMQTransport_GetMessage(t *testing.T) {
	transport := NewRMQTransport()
	err := transport.Init()
	if err != nil {
		t.Fatalf("error initializing RMQTransport: %v", err)
	}
	defer func() {
		err = transport.Destruct()
		if err != nil {
			t.Fatalf("error destructing RMQTransport: %v", err)
		}
	}()

	msg := &protocol.Message{
		Payload: []byte("test message"),
		Metadata: &protocol.Metadata{
			Sender:    "test-sender",
			Recipient: "test-recipient",
		},
	}

	err = transport.SendMessage(msg)
	if err != nil {
		t.Fatalf("error sending message: %v", err)
	}

	receivedMsg := transport.GetMessage()
	if receivedMsg == nil {
		t.Fatalf("expected to receive message, but got nil")
	}
}

func TestRMQTransport_AckMessage(t *testing.T) {
	transport := NewRMQTransport()
	err := transport.Init()
	if err != nil {
		t.Fatalf("error initializing RMQTransport: %v", err)
	}
	defer func() {
		err = transport.Destruct()
		if err != nil {
			t.Fatalf("error destructing RMQTransport: %v", err)
		}
	}()

	msg := &protocol.Message{
		Payload: []byte("test message"),
		Metadata: &protocol.Metadata{
			Sender:    "test-sender",
			Recipient: "test-recipient",
		},
	}

	err = transport.SendMessage(msg)
	if err != nil {
		t.Fatalf("error sending message: %v", err)
	}

	receivedMsg := transport.GetMessage()
	if receivedMsg == nil {
		t.Fatalf("expected to receive message, but got nil")
	}

	transport.AckMessage(receivedMsg)
}
