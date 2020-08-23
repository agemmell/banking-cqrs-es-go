package CheckingAccountService

import (
	"errors"
	"fmt"
	"reflect"
)

type Command interface {
	isCommand()
}

type StoresEvents interface {
	GetAllEvents() []interface{}
	PersistEvents(events ...interface{})
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
		err := account.OpenAccount(commandType.AccountID, commandType.Name)
		if err != nil {
			return err
		}
		cas.PersistEvents(account.GetNewEvents()...)

	default:
		commandStruct := reflect.TypeOf(commandType).String()
		return errors.New(fmt.Sprintf("unknown command %s", commandStruct))
	}

	return nil
}

func (cas *CheckingAccountService) PersistEvents(events ...Event) {
	// convert to interface{} for the event store to use
	e := make([]interface{}, len(events))
	for i, v := range events {
		e[i] = v
	}
	cas.eventStore.PersistEvents(e...)
}
