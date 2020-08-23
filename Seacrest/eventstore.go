package Seacrest

type EventStore struct {
	events []interface{}
}

func NewEventStore(events ...interface{}) *EventStore {
	return &EventStore{events}
}

func (es *EventStore) GetAllEvents() []interface{} {
	return es.events
}

func (es *EventStore) PersistEvents(events ...interface{}) {
	es.events = append(es.events, events...)
}
