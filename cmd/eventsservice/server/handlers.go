package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ikehakinyemi/ventickets/cmd/lib/msgqueue"
	"github.com/ikehakinyemi/ventickets/cmd/lib/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
)

type eventServiceHandler struct {
	dbhandler persistence.DatabaseHandler
	eventEmitter msgqueue.EventEmitter
}

func NewEventHandler(databasehandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) *eventServiceHandler {
	return &eventServiceHandler{
		dbhandler: databasehandler,
		eventEmitter: eventEmitter,
	}
}

func (eh *eventServiceHandler) FindEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "No search criteria found, you can either
						search by id via /id/4
						to search by name via /name/coldplayconcert}"`)
		return
	}

	searchkey, ok := vars["search"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "No search keys found, you can either search
						by id via /id/4
						to search by name via /name/coldplayconcert}"`)
		return
	}

	var event *persistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchkey)
	case "id":
		event, err = eh.dbhandler.FindEvent(searchkey)
	}
	if err != nil {
		fmt.Fprintf(w, `{"error": "%s"}`, err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) AllEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllAvailableEvents()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "Error occured while trying to find all available events %s"}`, err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "Error occured while trying encode events to JSON %s"}`, err)
	}
}

func (eh *eventServiceHandler) NewEventHandler(w http.ResponseWriter, r *http.Request) {
	event := persistence.Event{}

	err := json.NewDecoder(r.Body).Decode(&event)
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "error occured while decoding event data: %s"}`, err)
		return
	}
	id, err := eh.dbhandler.AddEvent(event)
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "error occured while persisting event: %s"}`, err)
		return
	}

	eventID, ok := id.(primitive.ObjectID)
	if !ok {
		w.WriteHeader(500)
		fmt.Fprintln(w, `{"error": "internal server error}`)
		return
	}

	fmt.Fprintf(w, `{"id":%q}`, eventID.Hex())
}

