package listener

import (
	"log"

	"github.com/ikehakinyemi/ventickets/cmd/contracts"
	"github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue"
	"github.com/ikehakinyemi/ventickets/cmd/lib/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventProcessor struct {
	EventListener msgqueue.EventListener
	Database persistence.DatabaseHandler
}

func (p *EventProcessor) ProcessEvents() error {
	log.Println("Listening for events...")

	received, errors, err := p.EventListener.Listen("event.created")
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-received:
			p.handleEvent(event)
		case err := <-errors:
			log.Printf("received an error while processing msg: %s", err.Error())
		}
	}
}

func (p *EventProcessor) handleEvent(event msgqueue.Event) {
	switch e := event.(type) {
	case *contracts.EventCreatedEvent:
		log.Printf("event %s created: %s", e.ID, e)
		id, err := primitive.ObjectIDFromHex(e.ID)
		if err != nil {
			log.Printf("event id %s is not a valid mongodb primitive.ObjectID: %s\n", e.ID, err)
			return
		}
		p.Database.AddEvent(persistence.Event{ID: id})
	case *contracts.LocationCreatedEvent:
		log.Printf("location %s created: %s", e.ID, e)
		id, err := primitive.ObjectIDFromHex(e.ID)
		if err != nil {
			log.Printf("location id %s is not a valid mongodb primitive.ObjectID: %s\n", e.ID, err)
			return
		}
		p.Database.AddLocation(persistence.Event{ID: id})
	default:
		log.Printf("unknown event: %t", e)
	}
}