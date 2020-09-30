package CheckingAccountService

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agemmell/banking-cqrs-es-go/Seacrest"
	uuid "github.com/nu7hatch/gouuid"
	"reflect"
	"time"
)

type Command interface {
	isCommand()
}

// StoreEvents: I'm directly coupling to my little event store which is not ideal but hey. It might be better to
//  have a repository-like struct to be the thing that gets coupled and does all the translating work between this
//  service's events and the event store's events. That way only the repo needs to be replaced should I change my event
//  store tech.
type StoresEvents interface {
	GetAllEvents() []Seacrest.EventEnvelope
	GetEventsByAggregateID(aggregateID string) map[uint]Seacrest.EventEnvelope
	WriteEventsToFile(filename string) error
	PersistEvent(aggregateID string, eventType string, payload []byte) error
	LoadEventsFromFile(filename string) error
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
		err = cas.PersistEvents(account.GetNewEvents()...)
		if err != nil {
			return err
		}

	case DepositMoney:
		events, err := cas.GetEventsByAggregateID(commandType.ID)
		if err != nil {
			return err
		}
		account := Account{}
		err = account.LoadFromEvents(events)
		if err != nil {
			return err
		}
		err = account.DepositMoney(commandType.Amount)
		if err != nil {
			return err
		}
		err = cas.PersistEvents(account.GetNewEvents()...)
		if err != nil {
			return err
		}

	case WithdrawMoney:
		events, err := cas.GetEventsByAggregateID(commandType.ID)
		if err != nil {
			return err
		}
		account := Account{}
		err = account.LoadFromEvents(events)
		if err != nil {
			return err
		}
		err = account.WithdrawMoney(commandType.Amount)
		if err != nil {
			return err
		}
		err = cas.PersistEvents(account.GetNewEvents()...)
		if err != nil {
			return err
		}

	case CloseAccount:
		events, err := cas.GetEventsByAggregateID(commandType.ID)
		if err != nil {
			return err
		}
		account := Account{}
		err = account.LoadFromEvents(events)
		if err != nil {
			return err
		}
		err = account.CloseAccount()
		if err != nil {
			return err
		}
		err = cas.PersistEvents(account.GetNewEvents()...)
		if err != nil {
			return err
		}

	default:
		commandStruct := reflect.TypeOf(commandType).String()
		return errors.New(fmt.Sprintf("unknown command %s", commandStruct))
	}

	return nil
}

func (cas *CheckingAccountService) PersistEvents(events ...Event) error {
	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}
		err = cas.eventStore.PersistEvent(event.AggregateID(), event.EventType(), payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cas *CheckingAccountService) GetAllEvents() ([]Event, error) {
	envelopes := cas.eventStore.GetAllEvents()
	var events []Event
	for _, envelope := range envelopes {
		event, err := cas.TransformEnvelopeToEvent(envelope)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (cas *CheckingAccountService) GetEventsByAggregateID(aggregateID string) ([]Event, error) {
	envelopes := cas.eventStore.GetEventsByAggregateID(aggregateID)
	var events []Event
	for _, envelope := range envelopes {
		event, err := cas.TransformEnvelopeToEvent(envelope)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// GenerateCheckingAccountEvents: for test purposes
func (cas *CheckingAccountService) GenerateCheckingAccountEvents() error {
	UUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	aggregateID := UUID.String()

	UUID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	aggregateID2 := UUID.String()

	var events []Event
	events = append(events,
		AccountWasOpened{
			ID:        aggregateID,
			Name:      "Alex Gemmell",
			Timestamp: time.Now().UnixNano(),
		}, MoneyWasDeposited{
			ID:        aggregateID,
			Amount:    12400,
			Timestamp: time.Now().UnixNano(),
		}, AccountWasOpened{
			ID:        aggregateID2,
			Name:      "Bobby Tables",
			Timestamp: time.Now().UnixNano(),
		}, MoneyWasDeposited{
			ID:        aggregateID,
			Amount:    1200,
			Timestamp: time.Now().UnixNano(),
		}, MoneyWasDeposited{
			ID:        aggregateID2,
			Amount:    144000,
			Timestamp: time.Now().UnixNano(),
		}, MoneyWasDeposited{
			ID:        aggregateID2,
			Amount:    1299,
			Timestamp: time.Now().UnixNano(),
		}, MoneyWasWithdrawn{
			ID:        aggregateID,
			Amount:    4200,
			Balance:   9400,
			Timestamp: time.Now().UnixNano(),
		}, MoneyWasWithdrawn{
			ID:        aggregateID2,
			Amount:    55288,
			Balance:   90011,
			Timestamp: time.Now().UnixNano(),
		}, MoneyWasDeposited{
			ID:     aggregateID,
			Amount: 999,
		},
	)

	err = cas.PersistEvents(events...)
	if err != nil {
		return err
	}

	return nil
}

func (cas *CheckingAccountService) WriteEventsToFile(filename string) error {
	err := cas.eventStore.WriteEventsToFile(filename)
	if err != nil {
		return err
	}

	return nil
}

func (cas *CheckingAccountService) LoadEventsFromFile(filename string) error {
	err := cas.eventStore.LoadEventsFromFile(filename)
	if err != nil {
		return err
	}

	return nil
}

func (cas *CheckingAccountService) TransformEnvelopeToEvent(envelope Seacrest.EventEnvelope) (Event, error) {
	var event Event
	switch envelope.EventType {
	case TypeAccountWasOpened:
		event = &AccountWasOpened{}
	case TypeMoneyWasDeposited:
		event = &MoneyWasDeposited{}
	case TypeMoneyWasWithdrawn:
		event = &MoneyWasWithdrawn{}
	case TypeWithdrawFailedDueToInsufficientFunds:
		event = &WithdrawFailedDueToInsufficientFunds{}
	case TypeAccountWasClosed:
		event = &AccountWasClosed{}
	default:
		return nil, errors.New(fmt.Sprintf("unknown event type in envelope %s", envelope.EventType))
	}

	err := json.Unmarshal(envelope.Payload, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (cas *CheckingAccountService) HydrateEvent(payload []byte, event Event) error {
	err := json.Unmarshal(payload, &event)
	if err != nil {
		return err
	}
	return nil
}
