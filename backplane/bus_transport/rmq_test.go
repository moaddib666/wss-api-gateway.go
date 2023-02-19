package bus_transport

import (
	"MargayGateway/protocol"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"testing"
	"time"
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

	publishTestMessage(t, msg)

	receivedMsg := transport.GetMessage()
	if receivedMsg == nil {
		t.Fatalf("expected to receive message, but got nil")
	}
}

func publishTestMessage(t *testing.T, msg *protocol.Message) {
	con, err := amqp.Dial(os.Getenv("MARGAY_TRANSPORT_DSN"))
	if err != nil {
		t.Fatalf("error connecting to RabbitMQ: %v", err)
	}
	ch, err := con.Channel()
	if err != nil {
		t.Fatalf("error creating channel: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = ch.PublishWithContext(
		ctx,
		inboxQueueName,
		"",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Body: msg.Payload,
			Headers: map[string]interface{}{
				"recipient": msg.Metadata.Recipient,
				"sender":    msg.Metadata.Sender,
			},
		})
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

	publishTestMessage(t, msg)

	receivedMsg := transport.GetMessage()
	if receivedMsg == nil {
		t.Fatalf("expected to receive message, but got nil")
	}

	transport.AckMessage(receivedMsg)
}
