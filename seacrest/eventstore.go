package seacrest

type MessageDescriber interface {
	MessageID() string
	MessageType() string
}

type Message struct {
	messageID   string
	messageType string
}

func (m *Message) MessageID() string {
	return m.messageID
}

func (m *Message) MessageType() string {
	return m.messageType
}

type EventStore struct {
	events []MessageDescriber
}

func NewEventStore(events ...MessageDescriber) *EventStore {
	return &EventStore{events}
}

func (es *EventStore) GetAllEvents() []Message {
	return es.events
}

func (es *EventStore) PersistEvents(events ...MessageDescriber) {
	es.events = append(es.events, events...)
}
