package main

import (
	"flag"
	"log"

	"github.com/Shopify/sarama"
	"github.com/ikehakinyemi/ventickets/cmd/bookingservices/listener"
	"github.com/ikehakinyemi/ventickets/cmd/bookingservices/rest"
	"github.com/ikehakinyemi/ventickets/cmd/lib/configuration"
	"github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue"
	msgqueue_amqp "github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue/amqp"
	"github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue/kafka"
	"github.com/ikehakinyemi/ventickets/cmd/lib/persistence/dblayer"
	"github.com/streadway/amqp"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var eventListener msgqueue.EventListener
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("conf", `./cmd/lib/configuration/config.json`, "flag to set the path to the configuration json file")
	connStrDB := flag.String("db-conn-str", "", "Set MongoDB connection string")
	flag.Parse()

	//extract configuration
	config, _ := configuration.ExtractConfiguration(*confPath, *connStrDB)

	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		panicIfErr(err)

		eventListener, err = msgqueue_amqp.NewAMQPEventListener(conn, "events", "booking")
		panicIfErr(err)

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		panicIfErr(err)
	case "kafka":
		conf := sarama.NewConfig()
		conf.Producer.Return.Successes = true
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		panicIfErr(err)

		eventListener, err = kafka.NewKafkaEventListener(conn, []int32{})
		panicIfErr(err)

		eventEmitter, err = kafka.NewKafkaEventEmitter(conn)
		panicIfErr(err)
	default:
		panic("Bad message broker type: " + config.MessageBrokerType)
	}

	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbhandler.Close()
	log.Println("Database connection successful")

	processor := listener.EventProcessor{EventListener: eventListener, Database: dbhandler}
	go processor.ProcessEvents()

	rest.ServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)
}
