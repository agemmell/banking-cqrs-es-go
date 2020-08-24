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
	id        string
	name      string
	balance   int
	version   uint
	newEvents []Event
	open      bool
}

func (a *Account) AggregateID() string {
	return a.id
}

func (a *Account) Version() uint {
	return a.version
}

func (a *Account) raiseEvent(event Event) error {
	err := a.ApplyEvent(event)
	if err != nil {
		return err
	}

	a.RecordNewEvent(event)
	return nil
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
		a.id = eventType.ID
		a.name = eventType.Name
		a.open = true
		a.balance = 0
		a.version = eventType.Version()
	case MoneyWasDeposited:
		a.balance += eventType.Amount
		a.version = eventType.Version()
	case MoneyWasWithdrawn:
		a.balance -= eventType.Amount
		a.version = eventType.Version()
	case WithdrawFailedDueToInsufficientFunds:
		a.version = eventType.Version()
	case AccountWasClosed:
		a.open = false
		a.version = eventType.Version()
	default:
		eventStruct := reflect.TypeOf(eventType).String()
		return errors.New(fmt.Sprintf("unknown event %s", eventStruct))
	}
	return nil
}

// LoadFromEvents: Return aggregate to state from past events without triggering side effects
func (a *Account) LoadFromEvents(events []Event) error {
	for _, event := range events {
		err := a.ApplyEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

// Command Handlers: protect aggregate invariants before throwing an event

// OpenAccount: open a new account
func (a *Account) OpenAccount(id string, name string) error {

	if len(a.id) > 0 && len(a.name) > 0 && a.version == 0 {
		return errors.New(fmt.Sprintf("cannot open an already open account [account: %+v]", a))
	}

	event := AccountWasOpened{
		ID:      id,
		Name:    name,
		version: a.version + 1,
	}

	return a.raiseEvent(event)
}

// DepositMoney: deposit money into an account
func (a *Account) DepositMoney(amount int) error {

	if a.open == false {
		return errors.New(fmt.Sprintf("cannot deposit money into an unopened account [account: %+v]", a))
	}

	if amount <= 0 {
		return errors.New(fmt.Sprintf("deposit amount must be greater than 1 [amount: %+v]", amount))
	}

	event := MoneyWasDeposited{
		ID:      a.id,
		Amount:  amount,
		version: a.version + 1,
	}

	return a.raiseEvent(event)
}

// WithdrawMoney: withdraw money from an account
func (a *Account) WithdrawMoney(amount int) error {

	if a.open == false {
		return errors.New(fmt.Sprintf("cannot withdraw money from an unopened account [account: %+v]", a))
	}

	if a.balance >= amount {
		event := MoneyWasWithdrawn{
			ID:      a.id,
			Amount:  amount,
			version: a.version + 1,
		}

		return a.raiseEvent(event)
	}

	event := WithdrawFailedDueToInsufficientFunds{
		ID:      a.id,
		Amount:  amount,
		version: a.version + 1,
	}

	return a.raiseEvent(event)
}

// CloseAccount: close the account
func (a *Account) CloseAccount() error {

	if a.open == false {
		return errors.New(fmt.Sprintf("cannot close a closed account [account: %+v]", a))
	}

	if a.balance > 0 {
		return errors.New(fmt.Sprintf("cannot close an account with a balance [account: %+v]", a))
	}

	event := AccountWasClosed{
		ID:      a.id,
		version: a.version + 1,
	}

	return a.raiseEvent(event)
}
