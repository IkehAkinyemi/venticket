package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue"
)

type KafkaEventEmitter struct {
	producer sarama.SyncProducer
}

func NewKafkaEventEmitter(client sarama.Client) (msgqueue.EventEmitter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	emitter := &KafkaEventEmitter{
		producer: producer,
	}

	return emitter, nil
}

func (e *KafkaEventEmitter) Emit(event msgqueue.Event) error {
	envelope := messageEnvelope{event.EventName(), event}
	json, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage {
		Topic: event.EventName(),
		Value: sarama.ByteEncoder(json),
	}

	_, _, err = e.producer.SendMessage(msg)
	return err
}
