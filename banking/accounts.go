package banking

import (
	"errors"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"

	"github.com/agemmell/banking-cqrs-es-go/seacrest"
)

// todo we're going to have account-related command handler functions in here

type AggregateDescriber interface {
	AggregateID() string
	Version() int
	ApplyEvent()
	GetEvents() []seacrest.MessageDescriber
}

type Aggregate struct {
	aggregateID string
	version     int
	events      []seacrest.MessageDescriber
}

type Account struct {
	accountID string
	name      string
	Aggregate
}

type AccountService struct {
	escqrs seacrest.Seacrest
}

func NewAccountService() *AccountService {
	return &AccountService{&seacrest.Service{}}
}

func NewAccount() *Account {
	account := Account{}
	account.version = 0
	return &account
}

func (account *Account) AggregateID() string {
	return account.aggregateID
}

func (account *Account) Version() int {
	return account.version
}

func (account *Account) ApplyEvent(event seacrest.MessageDescriber) error {
	switch eventType := event.(type) {
	case *AccountWasOpened:
		account.accountID = eventType.AccountID()
		account.name = eventType.Name()
		account.aggregateID = eventType.AccountID()
		account.version++
	default:
		return errors.New(fmt.Sprintf("unknown event type %v", eventType))
	}
	return nil
}

func (account *Account) GetEvents() []seacrest.MessageDescriber {
	return account.events
}


func (account *Account) RecordEvent(event seacrest.MessageDescriber) {
	account.events = append(account.events, event)
}

func (as *AccountService) CreateOpenAccount(accountID string, name string) (*OpenAccount, error) {
	message, err := as.escqrs.CreateMessageOfType(OpenAccountMessageType)
	if err != nil {
		return nil, err
	}
	return &OpenAccount{
		message,
		accountID,
		name,
	}, nil
}

func (as *AccountService) CreateAccountWasOpened(accountID string, name string) (*AccountWasOpened, error) {
	message, err := as.escqrs.CreateMessageOfType(AccountWasOpenedMessageType)
	if err != nil {
		return nil, err
	}
	return &AccountWasOpened{
		message,
		accountID,
		name,
	}, nil
}

func (as *AccountService) CreateOpenAccountWithUUID(accountID *uuid.UUID, name string) (*OpenAccount, error) {
	openAccount, err := as.CreateOpenAccount(accountID.String(), name)
	if err != nil {
		return nil, err
	}

	return openAccount, nil
}

func (as *AccountService) OpenAccount(command *OpenAccount) (*Account, error) {
	// No Account invariants to protect at this stage because Account doesn't exist yet
	account := NewAccount()
	accountWasOpened, err := as.CreateAccountWasOpened(command.AccountID(), command.Name())
	if err != nil {
		return nil, err
	}

	account.RecordEvent(accountWasOpened)

	err = account.ApplyEvent(accountWasOpened)
	if err != nil {
		return nil, err
	}
	return account, nil
}
