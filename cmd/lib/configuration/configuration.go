package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ikehakinyemi/ventickets/cmd/lib/persistence/dblayer"
)

var (
	DBTypeDefault       = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://127.0.0.1"
	RestfulEPDefault    = "localhost:8181"
	RestfulTLSEPDefault = "localhost:9191"
	AMQPMessageBroker = "amqp://guest:guest@localhost:5672"
	MessageBrokerType = "kafka"
	KafkaMessageBrokers = []string{"http://localhost:9092/"}
)

type ServiceConfig struct {
	Databasetype      dblayer.DBTYPE `json:"databasetype"`
	DBConnection      string         `json:"dbconnection"`
	RestfulEndpoint   string         `json:"restfulapi_endpoint"`
	RestfulTLSEndPint string         `json:"restfulapi-tlsendpoint"`
	AMQPMessageBroker string `json:"amqp_message_broker"`
	MessageBrokerType string `json:"message_broker_type"`
	KafkaMessageBrokers []string `json:"kafka_message_brokers"`
}

func ExtractConfiguration(filename string, connStr string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
		AMQPMessageBroker,
		MessageBrokerType,
		KafkaMessageBrokers,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}
	defer file.Close()

	if connStr != "" {
		conf.DBConnection = connStr
	}

	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		log.Fatal(err)
	}

	if broker := os.Getenv("AMQP_URL"); broker != "" {
		conf.AMQPMessageBroker = broker
	}

	if connStr != "" {
		conf.DBConnection = connStr
	}

	return conf, err
}
