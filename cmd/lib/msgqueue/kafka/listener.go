package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/ikehakinyemi/ventickets/cmd/contracts"
	"github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue"
	"github.com/mitchellh/mapstructure"
)

type KafkaEventListener struct {
	consumer sarama.Consumer
	partitions []int32
}

func NewKafkaEventListener (client sarama.Client, partitions []int32) (msgqueue.EventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	listener := &KafkaEventListener{
		consumer: consumer,
		partitions: partitions,
	}

	return listener, nil
}

func (l *KafkaEventListener) Listen(events ...string) (<-chan msgqueue.Event, <-chan error, error) {
	var err error
	topic := "events"
	results := make(chan msgqueue.Event)
	errors := make(chan error)

	partitions := l.partitions
	if len(partitions) == 0 {
		partitions, err = l.consumer.Partitions(topic)
		if err != nil {
			return nil, nil, err
		}
	}

	log.Printf("topic %s has %d partitions", topic, len(partitions))

	for _, partition := range partitions {
		conn, err := l.consumer.ConsumePartition(topic, partition, 0)
		if err != nil {
			return nil, nil, err
		}

		go func ()  {
			for msg := range conn.Messages() {
				body := messageEnvelope{}
				err = json.Unmarshal(msg.Value, &body)
				if err != nil {
					errors <- fmt.Errorf("could not JSON-decode message: %s", err)
					continue
				}

				var event msgqueue.Event
				switch body.EventName {
				case "event.created":
					event = &contracts.EventCreatedEvent{}
				case "location.created":
					event = &contracts.LocationCreatedEvent{}
				default:
					errors <- fmt.Errorf("unknown event type: %s", body.EventName)
					continue
				}

				cfg := mapstructure.DecoderConfig {
					Result: event,
					TagName: "json",
				}
				 decoder, err := mapstructure.NewDecoder(&cfg)
				 if err != nil {
					errors <- fmt.Errorf("could not map events %s: %s", body.EventName, err)
				 }
				decoder.Decode(body.Payload)
				results <- event
			}
		}()
	}

	return results, errors, nil
}