package Seacrest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewEventStore_GetAllEvents(t *testing.T) {
	eventStore := NewEventStore()
	assert.Len(t, eventStore.GetAllEvents(), 0)
}

type Event struct{}

func Test_EventStore_PersistEvents(t *testing.T) {
	t.Parallel()

	// Given
	event1 := Event{}
	event2 := Event{}
	event3 := Event{}
	eventStore := NewEventStore()

	// When
	eventStore.PersistEvents(event1, event2, event3)

	// Then
	assert.Len(t, eventStore.events, 3)
}
