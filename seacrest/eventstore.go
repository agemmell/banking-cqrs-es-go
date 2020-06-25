package seacrest

type StoresEvents interface {
	GetAllEvents() []MessageDescriber
	PersistEvents(events ...MessageDescriber)
}

type EventStore struct {
	events []MessageDescriber
}

func NewEventStore(events ...MessageDescriber) *EventStore {
	return &EventStore{events}
}

func (es *EventStore) GetAllEvents() []MessageDescriber {
	return es.events
}

func (es *EventStore) PersistEvents(events ...MessageDescriber) {
	es.events = append(es.events, events...)
}
