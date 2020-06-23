package cqrs_es

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewEventStore_GetAllEvents(t *testing.T) {
	eventStore := NewEventStore()
	assert.Len(t, eventStore.GetAllEvents(), 0)
}
