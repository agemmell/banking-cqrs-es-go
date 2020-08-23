package CheckingAccountService

import (
	"errors"
	"fmt"
	"reflect"
)

type Event interface {
	isEvent()
}

type Account struct {
	accountID string
	name      string
	version   int
	newEvents []Event
}

func (a *Account) AggregateID() string {
	return a.accountID
}

func (a *Account) Version() int {
	return a.version
}

func (a *Account) RecordNewEvent(event Event) {
	a.newEvents = append(a.newEvents, event)
}

func (a *Account) GetNewEvents() []Event {
	return a.newEvents
}

// ApplyEvent: Change aggregate state according to event type
func (a *Account) ApplyEvent(event Event) error {

	switch eventType := event.(type) {
	case AccountWasOpened:
		a.accountID = eventType.AccountID
		a.name = eventType.Name
		a.version = 1
	default:
		eventStruct := reflect.TypeOf(eventType).String()
		return errors.New(fmt.Sprintf("unknown event %s", eventStruct))
	}
	return nil
}

// LoadFromEvents: Return aggregate to state from past events without triggering side effects
func (a *Account) LoadFromEvents(events []Event) error {
	for _, event := range events {
		return a.ApplyEvent(event)
	}
	return nil
}

// Command Handlers: protect aggregate invariants before throwing an event

// OpenAccount: open a new account
func (a *Account) OpenAccount(accountID string, name string) error {

	if len(a.accountID) > 0 && len(a.name) > 0 && a.version == 0 {
		return errors.New(fmt.Sprintf("cannot open an already open account [account: %+v]", a))
	}

	event := AccountWasOpened{
		accountID,
		name,
	}

	err := a.ApplyEvent(event)
	if err != nil {
		return err
	}

	a.RecordNewEvent(event)

	return nil
}
