package mongolayer

import (
	"context"
	"time"

	"github.com/ikehakinyemi/ventickets/cmd/lib/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DB     = "venticket"
	USERS  = "users"
	EVENTS = "events"
)

func (dbLayer *MongoDBLayer) AddEvent(e persistence.Event) (interface{}, error) {
	collection := dbLayer.client.Database(DB).Collection(EVENTS)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	e.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, e)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, err
}

func (dbLayer *MongoDBLayer) FindEvent(id string) (*persistence.Event, error) {
	var event persistence.Event
	filter := bson.D{{Key: "_id", Value: id}}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	collection := dbLayer.client.Database(DB).Collection(EVENTS)
	err := collection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, err
}

func (dbLayer *MongoDBLayer) FindEventByName(name string) (*persistence.Event, error) {
	var event persistence.Event
	filter := bson.D{{Key: "name", Value: name}}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	collection := dbLayer.client.Database(DB).Collection(EVENTS)
	err := collection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, err
}

func (dbLayer *MongoDBLayer) FindAllAvailableEvents() ([]*persistence.Event, error) {
	var events []*persistence.Event
	firstCtx, firstCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer firstCancel()

	collection := dbLayer.client.Database(DB).Collection(EVENTS)
	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(firstCtx, bson.D{{}}, nil)
	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	secondCtx, secondCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer secondCancel()
	for cur.Next(secondCtx) {

		// create a value into which the single document can be decoded
		var elem persistence.Event
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		events = append(events, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	thirdCtx, thirdCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer thirdCancel()
	cur.Close(thirdCtx)

	return events, err

}
