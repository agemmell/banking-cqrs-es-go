package banking

import (
	uuid "github.com/nu7hatch/gouuid"

	cqrses "github.com/agemmell/banking-cqrs-es-go/cqrs-es"
)

type AccountService struct {
	cqrsesService cqrses.CQRSES
}

func (as *AccountService) CreateOpenAccount(accountID string, name string) (*OpenAccount, error) {
	message, err := as.cqrsesService.CreateMessageOfType(OpenAccountMessageType)
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
	message, err := as.cqrsesService.CreateMessageOfType(AccountWasOpenedMessageType)
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
