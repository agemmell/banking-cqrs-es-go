package main

import (
	"testing"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"

	"github.com/agemmell/banking-cqrs-es-go/banking"
	"github.com/agemmell/banking-cqrs-es-go/seacrest"
)

func Test_OpenAnAccount(t *testing.T) {
	// Given a Banking service
	eventStore := seacrest.NewEventStore()
	bankingService := banking.NewService(eventStore)

	// and an OpenAccount command
	accountService := banking.NewAccountService()
	accountID, err := uuid.NewV4()
	assert.Nil(t, err)
	openAccount, err := accountService.CreateOpenAccountWithUUID(accountID, "Alex Gemmell")
	assert.Nil(t, err)

	// When an OpenAccount command is sent to the Banking service
	bankingService.HandleCommand(openAccount)

	// Then an AccountWasOpened event is produced
	events := eventStore.GetAllEvents()
	assert.Len(t, events, 1)

	expectedAccountWasOpened, err := accountService.CreateAccountWasOpened(accountID.String(), "Alex Gemmell")
	assert.Nil(t, err)
	assert.Equal(t, expectedAccountWasOpened, events[0])
}

// // Move this to a banking service?
// func HandleCommand(account *accounts.OpenAccount) {
// 	// make the argument a generic command type
//
// }
