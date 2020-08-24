package Seacrest

type Event interface {
	AggregateID() string
	Version() uint
}

type EventStore struct {
	orderedEvents []Event                   // <global order> -> Event
	eventsByID    map[string]map[uint]Event // <aggregateID> -> <version> -> Event
}

func NewEventStore(events ...Event) *EventStore {
	eventsByID := make(map[string]map[uint]Event, 0)
	orderedEvents := make([]Event, 0)
	es := EventStore{orderedEvents, eventsByID}
	es.PersistEvents(events...)
	return &es
}

func (es *EventStore) GetAllEvents() []Event {
	return es.orderedEvents
}

func (es *EventStore) PersistEvents(events ...Event) {
	for _, e := range events {
		if _, ok := es.eventsByID[e.AggregateID()]; ok {
			if es.eventsByID[e.AggregateID()][e.Version()-1] != nil {
				continue
			}
			es.eventsByID[e.AggregateID()][e.Version()-1] = e
			es.orderedEvents = append(es.orderedEvents, e)
			continue
		}
		event := make(map[uint]Event, 0)
		event[e.Version()-1] = e
		es.eventsByID[e.AggregateID()] = event
		es.orderedEvents = append(es.orderedEvents, e)
	}
}

func (es *EventStore) GetEventsByAggregateID(aggregateID string) map[uint]Event {
	if aggregateEvents, ok := es.eventsByID[aggregateID]; ok {
		return aggregateEvents
	}
	return nil
}
