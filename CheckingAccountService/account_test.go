package CheckingAccountService

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type UnknownEvent struct{}

func (uc UnknownEvent) isEvent() {}

func TestAccount_ApplyEventUnknownType(t *testing.T) {
	t.Parallel()

	// Given
	account := Account{}
	unknownEvent := UnknownEvent{}

	// When
	err := account.ApplyEvent(unknownEvent)

	// Then
	assert.Equal(t, "unknown event CheckingAccountService.UnknownEvent", err.Error())
}

func TestAccount_RecordAndGetNewEvent(t *testing.T) {
	t.Parallel()

	// Given
	account := Account{}
	accountWasOpened := AccountWasOpened{
		AccountID: "ABCD",
		Name:      "Alex Gemmell",
	}

	// When
	account.RecordNewEvent(accountWasOpened)

	// Then
	newEvents := account.GetNewEvents()
	assert.Len(t, newEvents, 1)
	assert.Equal(t, accountWasOpened, newEvents[0])
}

func TestAccount_LoadFromEvents(t *testing.T) {
	t.Parallel()

	// Given
	accountWasOpened := AccountWasOpened{
		AccountID: "ABCD",
		Name:      "Alex Gemmell",
	}
	events := append([]Event{}, accountWasOpened)
	account := Account{}

	// When
	err := account.LoadFromEvents(events)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, accountWasOpened.AccountID, account.accountID)
	assert.Equal(t, accountWasOpened.Name, account.name)
	assert.Equal(t, 1, account.version)
	assert.Empty(t, account.newEvents)
}

func TestAccount_OpenAccount(t *testing.T) {
	t.Parallel()

	// Given
	accountID := "ABCD"
	name := "Alex Gemmell"
	account := Account{}

	// When
	err := account.OpenAccount(accountID, name)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, accountID, account.accountID)
	assert.Equal(t, name, account.name)
	assert.Equal(t, 1, account.version)
	assert.Len(t, account.newEvents, 1)
	assert.IsType(t, AccountWasOpened{}, account.newEvents[0])
}
