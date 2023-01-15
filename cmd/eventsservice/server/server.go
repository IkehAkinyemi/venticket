package server

import (
	"net/http"

	"github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue"
	"github.com/ikehakinyemi/ventickets/cmd/lib/persistence"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ServeAPI(endpoint, tlsendpoint string, databasehandler persistence.DatabaseHandler, emitter msgqueue.EventEmitter) (chan error, chan error) {
	handler := NewEventHandler(databasehandler, emitter)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.FindEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)
	httpErrChan := make(chan error)
	httptlsErrChan := make(chan error)

	server := handlers.CORS()(r)

	go func() {
		httptlsErrChan <- http.ListenAndServeTLS(tlsendpoint, "./cmd/eventsservice/cert.pem", "./cmd/eventsservice/key.pem", server)
	}()

	return httpErrChan, httptlsErrChan
}
