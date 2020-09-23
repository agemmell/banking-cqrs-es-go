package Seacrest

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewEventStore_GetAllEvents(t *testing.T) {
	eventStore := NewEventStore()
	assert.Len(t, eventStore.GetAllEvents(), 0)
}

type Event struct {
	id        string
	eventType string
	version   uint
	someValue string
}

func Test_EventStore_PersistEvent(t *testing.T) {
	t.Parallel()

	// Given
	event1 := Event{id: "A", eventType: "event1", someValue: "first"}
	event1Payload, err := json.Marshal(event1)
	assert.Nil(t, err)
	event2 := Event{id: "B", eventType: "event2", someValue: "second"}
	event2Payload, err := json.Marshal(event2)
	assert.Nil(t, err)
	event3 := Event{id: "C", eventType: "event3", someValue: "third"}
	event3Payload, err := json.Marshal(event3)
	assert.Nil(t, err)
	eventStore := NewEventStore()

	// When
	err = eventStore.PersistEvent(event1.id, event1.eventType, event1Payload)
	assert.Nil(t, err)
	err = eventStore.PersistEvent(event2.id, event2.eventType, event2Payload)
	assert.Nil(t, err)
	err = eventStore.PersistEvent(event3.id, event3.eventType, event3Payload)
	assert.Nil(t, err)

	// Then
	assert.Len(t, eventStore.orderedEvents, 3)
}

func Test_EventStore_GetEventsByAggregateID(t *testing.T) {
	t.Parallel()

	// Given
	event1 := Event{id: "A", eventType: "event1", someValue: "first"}
	event1Payload, err := json.Marshal(event1)
	assert.Nil(t, err)
	event2 := Event{id: "B", eventType: "event2", someValue: "second"}
	event2Payload, err := json.Marshal(event2)
	assert.Nil(t, err)
	event3 := Event{id: "C", eventType: "event3", someValue: "third"}
	event3Payload, err := json.Marshal(event3)
	assert.Nil(t, err)
	event4 := Event{id: "B", version: 2, eventType: "event4", someValue: "fourth"}
	event4Payload, err := json.Marshal(event4)
	assert.Nil(t, err)
	event5 := Event{id: "A", version: 2, eventType: "event5", someValue: "fifth"}
	event5Payload, err := json.Marshal(event5)
	assert.Nil(t, err)
	event6 := Event{id: "A", version: 3, eventType: "event6", someValue: "sixth"}
	event6Payload, err := json.Marshal(event6)
	assert.Nil(t, err)
	eventStore := NewEventStore()

	err = eventStore.PersistEvent(event1.id, event1.eventType, event1Payload)
	assert.Nil(t, err)
	err = eventStore.PersistEvent(event2.id, event2.eventType, event2Payload)
	assert.Nil(t, err)
	err = eventStore.PersistEvent(event3.id, event3.eventType, event3Payload)
	assert.Nil(t, err)
	err = eventStore.PersistEvent(event4.id, event4.eventType, event4Payload)
	assert.Nil(t, err)
	err = eventStore.PersistEvent(event5.id, event5.eventType, event5Payload)
	assert.Nil(t, err)
	err = eventStore.PersistEvent(event6.id, event6.eventType, event6Payload)
	assert.Nil(t, err)

	// When
	envelopes := eventStore.GetEventsByAggregateID("A")

	// Then
	assert.Len(t, envelopes, 3)
	assert.Equal(t, event1.id, envelopes[0].AggregateID)
	assert.Equal(t, event1.eventType, envelopes[0].EventType)
	assert.Equal(t, event5.id, envelopes[1].AggregateID)
	assert.Equal(t, event5.eventType, envelopes[1].EventType)
	assert.Equal(t, event6.id, envelopes[2].AggregateID)
	assert.Equal(t, event6.eventType, envelopes[2].EventType)
}
