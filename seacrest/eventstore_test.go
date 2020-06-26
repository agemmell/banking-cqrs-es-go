package seacrest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewEventStore_GetAllEvents(t *testing.T) {
	eventStore := NewEventStore()
	assert.Len(t, eventStore.GetAllEvents(), 0)
}

func Test_EventStore_PersistEvents(t *testing.T) {
	event1 := Message{
		messageID:   "test-id-1",
		messageType: "test-type-a",
	}

	event2 := Message{
		messageID:   "test-id-2",
		messageType: "test-type-b",
	}

	event3 := Message{
		messageID:   "test-id-3",
		messageType: "test-type-c",
	}

	eventStore := NewEventStore()
	eventStore.PersistEvents(event1, event2, event3)
	assert.Len(t, eventStore.events, 3)
}
