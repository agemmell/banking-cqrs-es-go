package CheckingAccountService

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type Event interface {
	AggregateID() string
	EventType() string
	EventTimestamp() int64
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
	case *AccountWasOpened:
		a.id = eventType.ID
		a.name = eventType.Name
		a.open = true
		a.balance = 0
	case *MoneyWasDeposited:
		a.balance += eventType.Amount
	case *MoneyWasWithdrawn:
		a.balance -= eventType.Amount
	case *WithdrawFailedDueToInsufficientFunds:
	case *AccountWasClosed:
		a.open = false
	default:
		eventStruct := reflect.TypeOf(eventType).String()
		return errors.New(fmt.Sprintf("unknown event %s", eventStruct))
	}
	a.version++
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
		ID:        id,
		Name:      name,
		Timestamp: time.Now().UnixNano(),
	}

	return a.raiseEvent(&event)
}

// DepositMoney: deposit money into an account
func (a *Account) DepositMoney(amount int) error {

	if a.open == false {
		return errors.New(fmt.Sprintf("cannot deposit money into an unopened account [account: %+v]", a))
	}

	if amount <= 0 {
		return errors.New(fmt.Sprintf("deposit Amount must be greater than 1 [Amount: %+v]", amount))
	}

	event := MoneyWasDeposited{
		ID:        a.id,
		Amount:    amount,
		Timestamp: time.Now().UnixNano(),
	}

	return a.raiseEvent(&event)
}

// WithdrawMoney: withdraw money from an account
func (a *Account) WithdrawMoney(amount int) error {

	if a.open == false {
		return errors.New(fmt.Sprintf("cannot withdraw money from an unopened account [account: %+v]", a))
	}

	if a.balance >= amount {
		event := MoneyWasWithdrawn{
			ID:        a.id,
			Amount:    amount,
			Balance:   a.balance - amount,
			Timestamp: time.Now().UnixNano(),
		}

		return a.raiseEvent(&event)
	}

	// TODO
	//  WithdrawFailedDueToInsufficientFunds isn't technically an event for the account aggregate because it doesn't change
	//  its state.
	//  A more advanced model might be to return an error here instead and have the account service command handler decide
	//  how to proceed.
	//  A) If commands are synchronous/blocking then we can return the necessary error to the user so they know the withdrawal failed.
	//  This is the preferable option, especially to start with, as it's more straight-forward to build and to reason with.
	//  E.g. Command comes in via HTTP, is processed successfully (e.g. state changing events generated), return a 200 OK
	//  (e.g. maybe with some data that contains a URL to retrieve the newly created or updated object)
	//  B) If commands are one-way ("asynchronous") then we return an ACK to the client to imply "message received and was queued"
	//  then the command is put on a Command Bus to be processed later/momentarily.  Prefer this for long running processes.
	//  E.g Command comes in via HTTP, is placed on a Command Bus successfully, return a 202 Accepted.
	//  References:
	//   - https://softwareengineering.stackexchange.com/questions/344052/how-to-handle-post-validation-errors-in-command-ddd-cqrs#344091
	//   - https://stackoverflow.com/questions/29916468/what-should-be-returned-from-the-api-for-cqrs-commands
	//   - https://groups.google.com/g/dddcqrs/c/pDHW7ErGNt0

	event := WithdrawFailedDueToInsufficientFunds{
		ID:        a.id,
		Amount:    amount,
		Balance:   a.balance,
		Timestamp: time.Now().UnixNano(),
	}

	return a.raiseEvent(&event)
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
		ID:        a.id,
		Timestamp: time.Now().UnixNano(),
	}

	return a.raiseEvent(&event)
}
