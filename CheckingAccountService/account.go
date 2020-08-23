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
	balance   int
	version   int
	newEvents []Event
	open      bool
}

func (a *Account) AggregateID() string {
	return a.accountID
}

func (a *Account) Version() int {
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
		a.accountID = eventType.AccountID
		a.name = eventType.Name
		a.open = true
		a.balance = 0
		a.version = 1
	case MoneyWasDeposited:
		a.balance += eventType.Amount
		a.version++
	case MoneyWasWithdrawn:
		a.balance -= eventType.Amount
		a.version++
	case WithdrawFailedDueToInsufficientFunds:
		a.version++
	case AccountWasClosed:
		a.open = false
		a.version++
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
func (a *Account) OpenAccount(accountID string, name string) error {

	if len(a.accountID) > 0 && len(a.name) > 0 && a.version == 0 {
		return errors.New(fmt.Sprintf("cannot open an already open account [account: %+v]", a))
	}

	event := AccountWasOpened{
		accountID,
		name,
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
		a.accountID,
		amount,
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
			a.accountID,
			amount,
		}

		return a.raiseEvent(event)
	}

	event := WithdrawFailedDueToInsufficientFunds{
		a.accountID,
		amount,
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
		a.accountID,
	}

	return a.raiseEvent(event)
}
