package bus_transport

import (
	"WSSFacade/protocol"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
)

const inboxQueueName = "ApiGatewayInbox"
const outboxQueueName = "ApiGatewayOutbox"
const connectionDSN = "amqp://user:bitnami@localhost:5672/" // TODO ENV VAR

type RMQTransport struct {
	conn    *amqp.Connection
	inbox   amqp.Queue
	msgChan chan *protocol.Message
	ackChan chan bool
}

func (s *RMQTransport) Init() (err error) {
	s.conn, err = amqp.Dial(s.getDSN())
	if err != nil {
		return err
	}
	err = s.createQueues()
	if err != nil {
		return err
	}
	go s.reader()

	return err
}

func (s *RMQTransport) createQueues() (err error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		inboxQueueName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	s.inbox, err = ch.QueueDeclare(
		inboxQueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = ch.QueueBind(
		inboxQueueName,
		"",
		inboxQueueName,
		false,
		nil)

	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		outboxQueueName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *RMQTransport) getDSN() string {
	dsn := os.Getenv("RMQ_DSN")
	if dsn == "" {
		dsn = connectionDSN
	}
	return dsn
}
func (s *RMQTransport) reader() error {
	ch, err := s.conn.Channel()
	if err != nil {
		log.Printf("Can't open stable channel to rmq %v", err)
		return err
	}
	log.Printf("starting event bus reader for queue `%s`", inboxQueueName)
	msgs, err := ch.Consume(
		inboxQueueName, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("can't consume messages from queue `%s` - %v", inboxQueueName, err)
		return err
	}
	for d := range msgs {
		log.Printf("Received a message: body:`%s` id:`%s` headers:`%v`", d.Body, d.MessageId, d.Headers)
		msg := msgConvertor(d)
		if msg == nil {
			log.Printf("Skip message: %s %s %v", d.Body, d.MessageId, d.Headers)
			d.Reject(false)
			continue
		}
		s.msgChan <- msg
		ack := <-s.ackChan
		if ack {
			_ = d.Ack(false)
		} else {
			log.Printf("Reschedule message: %s %s %v", d.Body, d.MessageId, d.Headers)
			time.Sleep(1 * time.Second)
			d.Reject(true)
		}

	}
	return nil
}

func (s *RMQTransport) Destruct() error {
	return s.conn.Close()
}

func (s *RMQTransport) GetMessage() *protocol.Message {
	msg := <-s.msgChan
	log.Printf("Found message in event bus for user `%s`", msg.Metadata.Recipient)
	return msg
}

func (s *RMQTransport) AckMessage(message *protocol.Message) {
	s.ackChan <- message != nil
}

func (s *RMQTransport) SendMessage(message *protocol.Message) error {
	outbox, err := s.conn.Channel()
	if err != nil {
		return err
	}
	log.Printf("Sending message `%s` to event bus - %+v", message.Payload, message.Metadata)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO make global context
	defer cancel()
	err = outbox.PublishWithContext(
		ctx,
		outboxQueueName,
		"",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			//ContentType: "text/plain",
			Body: message.Payload,
			//UserId:      message.Metadata.Sender,
			//AppId:       "ApiGateway",
			//Headers:
			Headers: map[string]interface{}{
				"recipient": message.Metadata.Recipient,
				"sender":    message.Metadata.Sender,
			},
		})
	return err
}

func NewRMQTransport() Transport {
	return &RMQTransport{
		msgChan: make(chan *protocol.Message),
		ackChan: make(chan bool),
	}
}

func msgConvertor(d amqp.Delivery) *protocol.Message {
	var sender, recipient string
	if d.Headers["sender"] == nil {
		sender = "unknown"
	} else {
		sender = d.Headers["sender"].(string)
	}
	if d.Headers["recipient"] == nil {
		return nil
	}
	recipient = d.Headers["recipient"].(string)
	return protocol.NewMessage(sender, recipient, d.Body)
}
