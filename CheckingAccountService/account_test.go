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
		ID:   "ABCD",
		Name: "Alex Gemmell",
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
		ID:   "ABCD",
		Name: "Alex Gemmell",
		version: 1,
	}
	events := append([]Event{}, accountWasOpened)
	account := Account{}

	// When
	err := account.LoadFromEvents(events)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, accountWasOpened.ID, account.id)
	assert.Equal(t, accountWasOpened.Name, account.name)
	assert.Equal(t, uint(1), account.version)
	assert.Empty(t, account.newEvents)
}

func TestAccount_OpenAccount(t *testing.T) {
	t.Parallel()

	// Given
	account := Account{}

	// When
	id := "ABCD"
	name := "Alex Gemmell"
	err := account.OpenAccount(id, name)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, id, account.id)
	assert.Equal(t, name, account.name)
	assert.Equal(t, 0, account.balance)
	assert.Equal(t, uint(1), account.version)
	assert.Len(t, account.newEvents, 1)
	assert.IsType(t, AccountWasOpened{}, account.newEvents[0])
}

func TestAccount_DepositMoney(t *testing.T) {
	t.Parallel()

	// Given
	id := "ABCD"
	name := "Alex Gemmell"
	accountWasOpened := AccountWasOpened{
		ID:   id,
		Name: name,
		version: 1,
	}
	events := append([]Event{}, accountWasOpened)
	account := Account{}
	err := account.LoadFromEvents(events)
	assert.Nil(t, err)

	// When
	amount := 1234

	err = account.DepositMoney(amount)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, id, account.id)
	assert.Equal(t, name, account.name)
	assert.Equal(t, amount, account.balance)
	assert.Equal(t, uint(2), account.version)
	assert.Len(t, account.newEvents, 1)
	assert.IsType(t, MoneyWasDeposited{}, account.newEvents[0])
}

func TestAccount_WithdrawMoney(t *testing.T) {
	t.Parallel()

	// Given
	id := "ABCD"
	name := "Alex Gemmell"
	accountWasOpened := AccountWasOpened{
		ID:   id,
		Name: name,
		version: 1,
	}
	depositAmount := 1099
	moneyWasDeposited := MoneyWasDeposited{
		ID:     id,
		Amount: depositAmount,
		version: 2,
	}
	events := append([]Event{}, accountWasOpened, moneyWasDeposited)
	account := Account{}
	err := account.LoadFromEvents(events)
	assert.Nil(t, err)

	// When
	withdrawAmount := 199

	err = account.WithdrawMoney(withdrawAmount)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, id, account.id)
	assert.Equal(t, name, account.name)
	assert.Equal(t, depositAmount-withdrawAmount, account.balance)
	assert.Equal(t, uint(3), account.version)
	assert.Len(t, account.newEvents, 1)
	assert.IsType(t, MoneyWasWithdrawn{}, account.newEvents[0])
}

func TestAccount_CloseAccount(t *testing.T) {
	t.Parallel()

	// Given
	id := "ABCD"
	name := "Alex Gemmell"
	accountWasOpened := AccountWasOpened{
		ID:   id,
		Name: name,
		version: 1,
	}
	depositAmount := 1099
	moneyWasDeposited := MoneyWasDeposited{
		ID:     id,
		Amount: depositAmount,
		version: 2,
	}
	withdrawAmount := 1099
	moneyWasWithdrawn := MoneyWasWithdrawn{
		ID:     id,
		Amount: withdrawAmount,
		version: 3,
	}
	events := append([]Event{}, accountWasOpened, moneyWasDeposited, moneyWasWithdrawn)
	account := Account{}
	err := account.LoadFromEvents(events)
	assert.Nil(t, err)

	// When
	err = account.CloseAccount()
	assert.Nil(t, err)

	// Then
	assert.Equal(t, id, account.id)
	assert.Equal(t, name, account.name)
	assert.Equal(t, 0, account.balance)
	assert.Equal(t, false, account.open)
	assert.Equal(t, uint(4), account.version)
	assert.Len(t, account.newEvents, 1)
	assert.IsType(t, AccountWasClosed{}, account.newEvents[0])
}
