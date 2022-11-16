package persistence

type DatabaseHandler interface {
	AddEvent(Event) ([]byte, error)
	FindEvent(string) (*Event, error)
	FindEventByName(string) (*Event, error)
	FindAllAvailableEvents() ([]*Event, error)
	Close()
}
