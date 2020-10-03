package CheckingAccountService

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agemmell/banking-cqrs-es-go/Seacrest"
	"github.com/icrowley/fake"
	uuid "github.com/nu7hatch/gouuid"
	"math/rand"
	"reflect"
	"sort"
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

// AutoGenerateEvents
func (cas *CheckingAccountService) AutoGenerateEvents() error {
	var allCustomerEvents []Event
	customerCount := 100000
	closeAccountChance := float32(0.1)
	randSeed := int64(99)
	rand.Seed(randSeed)
	fake.Seed(randSeed)

	for i := 0; i < customerCount; i++ {
		closeAccount := false
		if rand.Float32() >= closeAccountChance {
			closeAccount = true
		}
		events, err := cas.GenerateFakeCustomerEvents(closeAccount)
		if err != nil {
			return err
		}
		allCustomerEvents = append(allCustomerEvents, events...)
	}

	// TODO Improve this by making the event envelope timestamps use a faked clock. Currently all events have faked
	//  timestamps but the envelopes don't so all events are in the order in which they were generated (customer
	//  events are grouped together despite having mixed event timestamps)
	sort.Slice(allCustomerEvents, func(i, j int) bool {
		return allCustomerEvents[i].EventTimestamp() < allCustomerEvents[j].EventTimestamp()
	})
	err := cas.PersistEvents(allCustomerEvents...)
	if err != nil {
		return err
	}
	return nil
}

func (cas *CheckingAccountService) GenerateFakeCustomerEvents(closeAccount bool) ([]Event, error) {
	var events []Event

	startTime := time.Date(2020, time.January, 0, 0, 0, 0, 0, time.UTC)
	month := rand.Intn(9)
	day := rand.Intn(31)
	hour := time.Hour * time.Duration(rand.Intn(24))
	minute := time.Minute * time.Duration(rand.Intn(60))
	second := time.Second * time.Duration(rand.Intn(60))
	startTime.AddDate(0, month, day).Add(hour).Add(minute).Add(second)

	// generate a UUID
	UUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	aggregateID := UUID.String()

	// generate a fake name
	fullName := fake.FullName()

	// create AccountWasOpened event
	events = append(events, AccountWasOpened{
		ID:        aggregateID,
		Name:      fullName,
		Timestamp: startTime.UnixNano(),
	})

	// create random account deposit & withdraw events
	depositMax := rand.Intn(100000)
	balance := 0
	var depositAmount, withdrawnAmount int
	depositMin := 500
	maxEventsCount := 10

	for {
		day := rand.Intn(4)
		hour := time.Hour * time.Duration(rand.Intn(24))
		minute := time.Minute * time.Duration(rand.Intn(60))
		second := time.Second * time.Duration(rand.Intn(60))
		startTime.AddDate(0, month, day).Add(hour).Add(minute).Add(second)

		if depositMax == 0 {
			if closeAccount {
				if balance != 0 {
					withdrawnAmount = balance
					balance = 0
					events = append(events, MoneyWasWithdrawn{
						ID:        aggregateID,
						Amount:    withdrawnAmount,
						Balance:   balance,
						Timestamp: startTime.UnixNano(),
					})
				}

				events = append(events, AccountWasClosed{
					ID:        aggregateID,
					Timestamp: startTime.Add(time.Millisecond * 10).UnixNano(),
				})
			}
			break
		}

		withdraw := false
		withdrawChance := float32(0.9)
		if rand.Float32() <= withdrawChance {
			withdraw = true
		}
		// Must deposit if balance is 0
		if balance == 0 {
			withdraw = false
		}

		if withdraw {
			withdrawnAmount = rand.Intn(balance)
			balance -= withdrawnAmount
			events = append(events, MoneyWasWithdrawn{
				ID:        aggregateID,
				Amount:    withdrawnAmount,
				Balance:   balance,
				Timestamp: startTime.UnixNano(),
			})
		} else {
			if depositMax <= depositMin {
				depositAmount = depositMax
			} else {
				// TODO Bug for future Alex: this can potentially generate a deposit amount of zero!
				depositAmount = rand.Intn(depositMax)
			}
			depositMax -= depositAmount
			balance += depositAmount
			events = append(events, MoneyWasDeposited{
				ID:        aggregateID,
				Amount:    depositAmount,
				Timestamp: startTime.UnixNano(),
			})
		}

		if len(events) >= maxEventsCount {
			break
		}

	}

	return events, nil
}
