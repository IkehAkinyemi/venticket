package main

import (
	"flag"
	"log"

	"github.com/ikehakinyemi/ventickets/cmd/eventsservice/rest"
	"github.com/ikehakinyemi/ventickets/cmd/lib/configuration"
	"github.com/ikehakinyemi/ventickets/cmd/lib/persistence/dblayer"
)

func main() {
	confPath := flag.String("conf", `./cmd/lib/configuration/config.json`, "flag to set the path to the configuration json file")
	connStrDB := flag.String("db-conn-str", "", "Set MongoDB connection string")
	flag.Parse()

	//extract configuration
	config, _ := configuration.ExtractConfiguration(*confPath, *connStrDB)

	log.Println("Connecting to database")
	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbhandler.Close()

	log.Println("Database connection successful... ")

	//RESTful API start
	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndPint, dbhandler)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error: ", err)
	}
}
