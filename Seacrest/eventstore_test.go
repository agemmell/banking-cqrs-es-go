package Seacrest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewEventStore_GetAllEvents(t *testing.T) {
	eventStore := NewEventStore()
	assert.Len(t, eventStore.GetAllEvents(), 0)
}

type eventMessage struct {
	id      string
	version uint
}

func (e eventMessage) AggregateID() string {
	return e.id
}

func (e eventMessage) Version() uint {
	return e.version
}

func Test_EventStore_PersistEvents(t *testing.T) {
	t.Parallel()

	// Given
	event1 := eventMessage{id: "A", version: 1}
	event2 := eventMessage{id: "B", version: 1}
	event3 := eventMessage{id: "C", version: 1}
	event4 := eventMessage{id: "B", version: 2}
	event5 := eventMessage{id: "A", version: 2}
	event6 := eventMessage{id: "A", version: 3}
	eventStore := NewEventStore()

	// When
	eventStore.PersistEvents(event1, event2, event3, event4, event5, event6)

	// Then
	assert.Len(t, eventStore.orderedEvents, 6)
}

func Test_EventStore_GetEventsByAggregateID(t *testing.T) {
	t.Parallel()

	// Given
	event1 := eventMessage{id: "A", version: 1}
	event2 := eventMessage{id: "B", version: 1}
	event3 := eventMessage{id: "C", version: 1}
	event4 := eventMessage{id: "B", version: 2}
	event5 := eventMessage{id: "A", version: 2}
	event6 := eventMessage{id: "A", version: 3}
	eventStore := NewEventStore()
	eventStore.PersistEvents(event1, event2, event3, event4, event5, event6)

	// When
	events := eventStore.GetEventsByAggregateID("A")

	// Then
	assert.Len(t, events, 3)
	assert.Equal(t, event1, events[0])
	assert.Equal(t, event5, events[1])
	assert.Equal(t, event6, events[2])
}
