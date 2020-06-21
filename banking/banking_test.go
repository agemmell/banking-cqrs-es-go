package banking

import (
	"testing"

	cqrses "github.com/agemmell/banking-cqrs-es-go/cqrs-es"
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

	event1 := AccountWasOpened{
		&FakeMessage{"1234ABCD", "FakeTypeOne"},
		"5678EFGH",
		"Alex Gemmell",
	}
	event2 := AccountWasOpened{
		&FakeMessage{"9012IJKL", "FakeTypeTwo"},
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

}
