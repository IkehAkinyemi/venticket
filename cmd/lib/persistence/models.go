package persistence

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	First    string
	Last     string
	Age      int
	Bookings []Booking
}

func (u *User) String() string {
	return fmt.Sprintf("id: %s, first_name: %s, last_name: %s, Age: %d, Bookings: %v", u.ID, u.First, u.Last, u.Age, u.Bookings)
}

type Booking struct {
	Date    int64  `json:"date,omitempty"`
	EventID []byte
	Seats   int    `json:"seats,omitempty"`
}

type Event struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name      string      `json:"name,omitempty"`
	Duration  int         `json:"duration,omitempty"`
	StartDate int64       `json:"start_date,omitempty"`
	EndDate   int64       `json:"end_date,omitempty"`
	Location  Location    `json:"location,omitempty"`
}

type Location struct {
	Name      string      `json:"name,omitempty"`
	Address   string      `json:"address,omitempty"`
	Country   string      `json:"country,omitempty"`
	OpenTime  int         `json:"open_time,omitempty"`
	CloseTime int         `json:"close_time,omitempty"`
	Halls     []Hall      `json:"halls,omitempty"`
}

type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}
