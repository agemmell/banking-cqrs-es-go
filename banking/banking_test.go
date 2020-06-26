package banking

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agemmell/banking-cqrs-es-go/seacrest"
	"github.com/agemmell/banking-cqrs-es-go/seacrest/seacrestfakes"
)

func Test_NewServiceNoEvents(t *testing.T) {
	t.Parallel()

	eventStore := seacrest.NewEventStore()
	got := NewService(eventStore)
	assert.Len(t, got.eventStore.GetAllEvents(), 0)
}

func Test_NewServiceWithEvents(t *testing.T) {
	t.Parallel()

	fakeMessage1 := seacrestfakes.FakeMessageDescriber{}
	fakeMessage1.MessageIDReturns("1234ABCD")
	fakeMessage1.MessageTypeReturns("FakeTypeOne")
	event1 := AccountWasOpened{
		&fakeMessage1,
		"5678EFGH",
		"Alex Gemmell",
	}

	fakeMessage2 := seacrestfakes.FakeMessageDescriber{}
	fakeMessage2.MessageIDReturns("9012IJKL")
	fakeMessage2.MessageTypeReturns("FakeTypeTwo")
	event2 := AccountWasOpened{
		&fakeMessage2,
		"3456MNOP",
		"Andrew Garfield",
	}

	eventStore := seacrest.NewEventStore(event1, event2)
	got := NewService(eventStore)
	allEvents := got.eventStore.GetAllEvents()
	assert.Len(t, allEvents, 2)
	assert.Equal(t, allEvents[0].MessageID(), event1.MessageID())
	assert.Equal(t, allEvents[1].MessageID(), event2.MessageID())
}

func TestService_HandleCommand(t *testing.T) {
	t.Parallel()

	// Given
	//eventStore := seacrest.NewEventStore()
	//got := NewService(eventStore)
	// todo
}

func Test_HandleCommand_UnknownMessage(t *testing.T) {
	t.Parallel()

	// Given
	eventStore := seacrest.NewEventStore()
	bankingService := NewService(eventStore)

	fakeMessageUnknownType := seacrestfakes.FakeMessageDescriber{}
	fakeMessageUnknownType.MessageTypeReturns("test-type")

	// When
	err := bankingService.HandleCommand(&fakeMessageUnknownType)

	// Then
	assert.Equal(t, "unknown command type test-type", err.Error())
}
