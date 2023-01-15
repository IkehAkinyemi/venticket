package main

import (
	"flag"
	"log"

	"github.com/ikehakinyemi/ventickets/cmd/eventsservice/server"
	"github.com/ikehakinyemi/ventickets/cmd/lib/configuration"
	msgqueue_amqp "github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue/amqp"
	"github.com/ikehakinyemi/ventickets/cmd/lib/persistence/dblayer"
	"github.com/streadway/amqp"
)

func main() {
	confPath := flag.String("conf", `./cmd/lib/configuration/config.json`, "flag to set the path to the configuration json file")
	connStrDB := flag.String("db-conn-str", "", "Set MongoDB connection string")
	flag.Parse()

	//extract configuration
	config, _ := configuration.ExtractConfiguration(*confPath, *connStrDB)

	log.Println("Connecting to AMQP Broker...")
	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//setup dependency injection for AMQP
	emitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to database...")
	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbhandler.Close()
	log.Println("Database connection successful")

	//RESTful API start
	httpErrChan, httptlsErrChan := server.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndPint, dbhandler, emitter)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error: ", err)
	}
}
