package CheckingAccountService

import (
	"github.com/agemmell/banking-cqrs-es-go/Seacrest"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewServiceNoEvents(t *testing.T) {
	t.Parallel()

	eventStore := Seacrest.NewEventStore()
	checkingAccountService := New(eventStore)
	assert.Len(t, checkingAccountService.GetAllEvents(), 0)
}

type UnknownCommand struct{}

func (uc UnknownCommand) isCommand() {}

func Test_Handle_Unknown_Command(t *testing.T) {
	t.Parallel()

	// Given
	eventStore := Seacrest.NewEventStore()
	checkingAccountService := New(eventStore)

	unknownCommand := UnknownCommand{}

	// When
	err := checkingAccountService.HandleCommand(unknownCommand)

	// Then
	assert.Equal(t, "unknown command CheckingAccountService.UnknownCommand", err.Error())
}

func Test_OpenAnAccount(t *testing.T) {
	t.Parallel()

	// Given
	eventStore := Seacrest.NewEventStore()
	checkingAccountService := New(eventStore)

	accountUUID, err := uuid.NewV4()
	assert.Nil(t, err)
	openAccount := OpenAccount{
		ID:   accountUUID.String(),
		Name: "Alex Gemmell",
	}

	// When
	err = checkingAccountService.HandleCommand(openAccount)
	assert.Nil(t, err)

	// Then
	events := checkingAccountService.GetAllEvents()
	assert.Len(t, events, 1)

	eventType, ok := events[0].(AccountWasOpened)
	assert.True(t, ok)
	assert.Equal(t, openAccount.ID, eventType.ID)
	assert.Equal(t, openAccount.Name, eventType.Name)
}
