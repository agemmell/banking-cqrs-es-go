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
	accountID, err := uuid.NewV4()
	assert.Nil(t, err)
	accountName := "Alex Gemmell"
	accountService := banking.NewAccountService()
	openAccount, err := accountService.CreateOpenAccountWithUUID(accountID, accountName)
	assert.Nil(t, err)

	// When an OpenAccount command is sent to the Banking service
	err = bankingService.HandleCommand(openAccount)
	assert.Nil(t, err)

	// Then an AccountWasOpened event is produced
	events := eventStore.GetAllEvents()
	assert.Len(t, events, 1)

	eventType, ok := events[0].(*banking.AccountWasOpened)
	assert.True(t, ok)
	assert.Equal(t, accountID.String(), eventType.AccountID())
	assert.Equal(t, accountName, eventType.Name())
}
