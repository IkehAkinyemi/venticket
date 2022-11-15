package amqp

import (
	"encoding/json"

	"github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue"
	"github.com/streadway/amqp"
)

type AMQPEventEmitter struct {
	connection *amqp.Connection
}

func NewAMQPEventEmitter(conn *amqp.Connection) (msgqueue.EventEmitter, error) {
	emitter := &AMQPEventEmitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return nil, err
	}

	return emitter, nil
}

func (a *AMQPEventEmitter) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)
}

func (a * AMQPEventEmitter) Emit(event msgqueue.Event) error {
	wireData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	msg := amqp.Publishing {
		Headers: amqp.Table{"x-event-name": event.EventName()},
		Body: wireData,
		ContentType: "application/json",
	}

	return channel.Publish(
		"events",
		event.EventName(),
		false,
		false,
		msg,
	)
}
