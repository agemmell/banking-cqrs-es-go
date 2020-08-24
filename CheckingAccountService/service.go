package CheckingAccountService

import (
	"errors"
	"fmt"
	"github.com/agemmell/banking-cqrs-es-go/Seacrest"
	"reflect"
)

type Command interface {
	isCommand()
}

type StoresEvents interface {
	GetAllEvents() []Seacrest.Event
	PersistEvents(events ...Seacrest.Event)
	GetEventsByAggregateID(aggregateID string) map[uint]Seacrest.Event
}

type CheckingAccountService struct {
	eventStore StoresEvents
}

func New(eventStore StoresEvents) CheckingAccountService {
	return CheckingAccountService{eventStore}
}

// HandleCommand: Handles commands
func (cas *CheckingAccountService) HandleCommand(command Command) error {

	switch commandType := command.(type) {
	case OpenAccount:
		account := Account{}
		err := account.OpenAccount(commandType.ID, commandType.Name)
		if err != nil {
			return err
		}
		cas.PersistEvents(account.GetNewEvents()...)

	case DepositMoney:
		events := cas.GetEventsByAggregateID(commandType.ID)
		account := Account{}
		err := account.LoadFromEvents(events)
		if err != nil {
			return err
		}
		err = account.DepositMoney(commandType.Amount)
		if err != nil {
			return err
		}
		cas.PersistEvents(account.GetNewEvents()...)

	default:
		commandStruct := reflect.TypeOf(commandType).String()
		return errors.New(fmt.Sprintf("unknown command %s", commandStruct))
	}

	return nil
}

func (cas *CheckingAccountService) GetAllEvents() []Event {
	events := cas.eventStore.GetAllEvents()
	e := make([]Event, len(events))
	for i, v := range events {
		e[i] = v.(Event)
	}
	return e
}

func (cas *CheckingAccountService) PersistEvents(events ...Event) {
	// convert to Seacrest.Event for the event store to use
	e := make([]Seacrest.Event, len(events))
	for i, v := range events {
		e[i] = v.(Seacrest.Event)
	}
	cas.eventStore.PersistEvents(e...)
}

func (cas *CheckingAccountService) GetEventsByAggregateID(aggregateID string) []Event {
	events := cas.eventStore.GetEventsByAggregateID(aggregateID)
	e := make([]Event, len(events))
	for i, v := range events {
		e[i] = v.(Event)
	}
	return e
}
