package mongolayer

import "github.com/ikehakinyemi/ventickets/cmd/lib/persistence"

const (
	LOCATIONS = "locations"
)

func (l *MongoDBLayer) AddLocation(e persistence.Event) {}

func (l *MongoDBLayer) AddBookingForUser(d []byte, booking persistence.Booking) {}