package banking

import (
	"testing"

	cqrses "github.com/agemmell/banking-cqrs-es-go/cqrs-es"
	cqrsesfakes "github.com/agemmell/banking-cqrs-es-go/cqrs-es/cqrs-esfakes"
)

func Test_NewServiceNoEvents(t *testing.T) {
	t.Parallel()

	eventStore := cqrses.NewEventStore()
	got := NewService(eventStore)
	if len(got.eventStore.GetAllEvents()) != 0 {
		t.Errorf("got %v, want %v", got.eventStore, eventStore)
	}
}

func Test_NewServiceWithEvents(t *testing.T) {
	t.Parallel()

	fakeMessage1 := cqrsesfakes.FakeMessageDescriber{}
	fakeMessage1.MessageIDReturns("1234ABCD")
	fakeMessage1.MessageTypeReturns("FakeTypeOne")
	event1 := AccountWasOpened{
		&fakeMessage1,
		"5678EFGH",
		"Alex Gemmell",
	}

	fakeMessage2 := cqrsesfakes.FakeMessageDescriber{}
	fakeMessage2.MessageIDReturns("9012IJKL")
	fakeMessage2.MessageTypeReturns("FakeTypeTwo")
	event2 := AccountWasOpened{
		&fakeMessage2,
		"3456MNOP",
		"Andrew Garfield",
	}

	eventStore := cqrses.NewEventStore(event1, event2)
	got := NewService(eventStore)
	allEvents := got.eventStore.GetAllEvents()
	if len(allEvents) != 2 {
		t.Errorf("got %v, want %v", got.eventStore, eventStore)
	}
	if allEvents[0].MessageID() != event1.MessageID() {
		t.Errorf("got %v, want %v", allEvents[0].MessageID(), event1.MessageID())
	}
	if allEvents[1].MessageType() != event2.MessageType() {
		t.Errorf("got %v, want %v", allEvents[1].MessageType(), event1.MessageType())
	}
}

func TestService_HandleCommand(t *testing.T) {
	t.Parallel()

	// todo
}
