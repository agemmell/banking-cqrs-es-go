package cqrs_es

type StoresEvents interface {
	GetAllEvents() []MessageDescriber
}

type EventStore struct {
	events []MessageDescriber
}

func NewEventStore(events ...MessageDescriber) EventStore {
	return EventStore{events}
}

func (es EventStore) GetAllEvents() []MessageDescriber {
	return es.events
}
