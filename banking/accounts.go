package banking

import (
	"errors"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
)

// todo we're going to have account-related command handler functions in here

type AggregateDescriber interface {
	AggregateID() string
	Version() int
	ApplyEvent()
	GetEvents() []MessageDescriber
}

type Aggregate struct {
	aggregateID string
	version     int
	events      []MessageDescriber
}

type Account struct {
	accountID string
	name      string
	Aggregate
}

type MessageDescriber interface {
	MessageID() string
	MessageType() string
}

type GeneratesUUIDs interface {
	CreateUUIDString() (string, error)
}

// What if this was unexported? accountService?
// Would the NewAccountService func still be able to return it outside of the package?
type AccountService struct {
	uuidService GeneratesUUIDs
}

func NewAccountService(uuidService GeneratesUUIDs) *AccountService {
	return &AccountService{uuidService}
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

func (account *Account) ApplyEvent(event MessageDescriber) error {
	switch eventType := event.(type) {
	case *AccountWasOpened:
		account.accountID = eventType.AccountID()
		account.name = eventType.Name()
		account.aggregateID = eventType.AccountID()
		account.version++
	default:
		return errors.New(fmt.Sprintf("unknown event type %v", eventType.MessageType()))
	}
	return nil
}

func (account *Account) GetEvents() []MessageDescriber {
	return account.events
}

func (account *Account) GetNewEvents() []MessageDescriber {
	/*
		todo How do we apply events and retrieve new events for persisting to event store?
		todo What about aggregate version number? Am I doing that right?
	*/
	return []MessageDescriber{}
}

func (account *Account) RecordEvent(event MessageDescriber) {
	account.events = append(account.events, event)
}

func (as *AccountService) CreateOpenAccount(accountID string, name string) (*OpenAccount, error) {
	messageID, err := as.uuidService.CreateUUIDString()
	if err != nil {
		return nil, err
	}
	return &OpenAccount{
		Message{
			messageID,
			OpenAccountMessageType,
		},
		accountID,
		name,
	}, nil
}

func (as *AccountService) CreateAccountWasOpened(accountID string, name string) (*AccountWasOpened, error) {
	messageID, err := as.uuidService.CreateUUIDString()
	if err != nil {
		return nil, err
	}
	return &AccountWasOpened{
		Message{
			messageID,
			AccountWasOpenedMessageType,
		},
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

//func (as *AccountService) CreateDepositMoney(accountID string, depositAmount int) (*DepositMoney, error) {
//	return DepositMoney
//}
