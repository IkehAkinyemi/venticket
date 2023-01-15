package contracts

import "time"

type EventCreatedEvent struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	LocationName string    `json:"location_name"`
	Start        time.Time `json:"start_time"`
	End          time.Time `json:"end_time"`
}

func (e *EventCreatedEvent) EventName() string {
	return "event.created"
}

func (e *EventCreatedEvent) PartitionerKey() string {
	return e.ID
}

type LocationCreatedEvent struct {
	ID string `json:"id"`
}

func (l *LocationCreatedEvent) EventName() string {
	return "location.created"
}

func (l *LocationCreatedEvent) PartitionerKey() string {
	return l.ID
}

type EventBookedEvent struct {
	EventID string
	UserID  string
}

func (l *EventBookedEvent) EventName() string {
	return "events.booked"
}

func (l *EventBookedEvent) PartitionerKey() string {
	return l.EventID
}