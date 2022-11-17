package persistence

type DatabaseHandler interface {
	AddEvent(Event) (interface{}, error)
	FindEvent(string) (*Event, error)
	FindEventByName(string) (*Event, error)
	FindAllAvailableEvents() ([]*Event, error)
	Close()
}
