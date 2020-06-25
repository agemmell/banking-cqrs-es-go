package banking

import (
	"errors"

	"github.com/agemmell/banking-cqrs-es-go/seacrest"
)

type Banking struct {
	accountService AccountService
	eventStore     seacrest.StoresEvents
}

func NewService(eventStore seacrest.StoresEvents) Banking {
	return Banking{NewAccountService(), eventStore}
}

func (b *Banking) HandleCommand(command seacrest.MessageDescriber) error {
	switch command.MessageType() {
	case OpenAccountMessageType:
		openAccount, ok := command.(OpenAccount)
		if ok != true {
			return errors.New("command has wrong message type")
		}

		err := b.accountService.OpenAccount(&openAccount)
		if err != nil {
			return err
		}
	}
	return nil
}
