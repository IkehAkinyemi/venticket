package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// uri := "mongodb+srv://" + url.QueryEscape("IkehAkinyemi") + ":" + 
	// 	url.QueryEscape("password") + "@" + "central.ox9gkqx.mongodb.net" + 
	// 	"/events?retryWrites=true&w=majority" "mongodb+srv://IkehAkinyemi:36C%23DNDjs9%21r%2A6a@cluster0.uw55o.mongodb.net/?retryWrites=true&w=majority"

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(	"mongodb+srv://IkehAkinyemi:36C%23DNDjs9%21r%2A6a@central.ox9gkqx.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}
