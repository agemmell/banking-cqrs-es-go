package banking

import (
	"errors"
	"fmt"
	"github.com/agemmell/banking-cqrs-es-go/seacrest"
)

type Banking struct {
	accountService *AccountService
	eventStore     seacrest.StoresEvents
}

func NewService(eventStore seacrest.StoresEvents) Banking {
	return Banking{NewAccountService(), eventStore}
}

func (b *Banking) HandleCommand(command seacrest.MessageDescriber) error {
	switch commandType := command.(type) {
	case *OpenAccount:
		account, err := b.accountService.OpenAccount(commandType)
		if err != nil {
			return err
		}
		b.eventStore.PersistEvents(account.GetEvents()...)

	default:
		return errors.New(fmt.Sprintf("unknown command type %v", commandType.MessageType()))
	}

	return nil
}
