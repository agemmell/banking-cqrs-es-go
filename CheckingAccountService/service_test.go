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
	allEvents, err := checkingAccountService.GetAllEvents()
	assert.Nil(t, err)
	assert.Len(t, allEvents, 0)
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

	// When
	accountUUID, err := uuid.NewV4()
	assert.Nil(t, err)
	openAccount := OpenAccount{
		ID:   accountUUID.String(),
		Name: "Alex Gemmell",
	}
	err = checkingAccountService.HandleCommand(openAccount)
	assert.Nil(t, err)

	// Then
	events, err := checkingAccountService.GetAllEvents()
	assert.Nil(t, err)
	assert.Len(t, events, 1)

	eventType, ok := events[0].(*AccountWasOpened)
	assert.True(t, ok)
	assert.Equal(t, openAccount.ID, eventType.ID)
	assert.Equal(t, openAccount.Name, eventType.Name)
}

func Test_DepositMoney(t *testing.T) {
	t.Parallel()

	// Given
	eventStore := Seacrest.NewEventStore()
	checkingAccountService := New(eventStore)
	id := "ABCD"
	name := "Alex Gemmell"
	accountWasOpened := AccountWasOpened{
		ID:   id,
		Name: name,
	}
	historicalEvents := append([]Event{}, accountWasOpened)
	err := checkingAccountService.PersistEvents(historicalEvents...)
	assert.Nil(t, err)

	// When
	depositMoney := DepositMoney{
		ID:   id,
		Amount: 1099,
	}
	err = checkingAccountService.HandleCommand(depositMoney)
	assert.Nil(t, err)

	// Then
	events, err := checkingAccountService.GetAllEvents()
	assert.Nil(t, err)
	assert.Len(t, events, 2)

	eventType, ok := events[1].(*MoneyWasDeposited)
	assert.True(t, ok)
	assert.Equal(t, depositMoney.ID, eventType.ID)
	assert.Equal(t, depositMoney.Amount, eventType.Amount)
}

func Test_WithdrawMoney(t *testing.T) {
	t.Parallel()

	// Given
	eventStore := Seacrest.NewEventStore()
	checkingAccountService := New(eventStore)
	id := "ABCD"
	name := "Alex Gemmell"
	accountWasOpened := AccountWasOpened{
		ID:   id,
		Name: name,
	}
	moneyWasDeposited := MoneyWasDeposited{
		ID:     id,
		Amount: 1099,
	}
	historicalEvents := append([]Event{}, accountWasOpened, moneyWasDeposited)
	err := checkingAccountService.PersistEvents(historicalEvents...)
	assert.Nil(t, err)

	// When
	withdrawMoney := WithdrawMoney{
		ID:   id,
		Amount: 199,
	}
	err = checkingAccountService.HandleCommand(withdrawMoney)
	assert.Nil(t, err)

	// Then
	events, err := checkingAccountService.GetAllEvents()
	assert.Nil(t, err)
	assert.Len(t, events, 3)

	eventType, ok := events[2].(*MoneyWasWithdrawn)
	assert.True(t, ok)
	assert.Equal(t, withdrawMoney.ID, eventType.ID)
	assert.Equal(t, withdrawMoney.Amount, eventType.Amount)
	assert.Equal(t, 900, eventType.Balance)
}

func Test_WithdrawFailedDueToInsufficientFunds(t *testing.T) {
	t.Parallel()

	// Given
	eventStore := Seacrest.NewEventStore()
	checkingAccountService := New(eventStore)
	id := "ABCD"
	name := "Alex Gemmell"
	accountWasOpened := AccountWasOpened{
		ID:   id,
		Name: name,
	}
	moneyWasDeposited := MoneyWasDeposited{
		ID:     id,
		Amount: 1099,
	}
	moneyWasWithdrawn := MoneyWasWithdrawn{
		ID:     id,
		Amount: 1099,
	}
	historicalEvents := append([]Event{}, accountWasOpened, moneyWasDeposited, moneyWasWithdrawn)
	err := checkingAccountService.PersistEvents(historicalEvents...)
	assert.Nil(t, err)

	// When
	withdrawMoney := WithdrawMoney{
		ID:   id,
		Amount: 1,
	}
	err = checkingAccountService.HandleCommand(withdrawMoney)
	assert.Nil(t, err)

	// Then
	events, err := checkingAccountService.GetAllEvents()
	assert.Nil(t, err)
	assert.Len(t, events, 4)

	eventType, ok := events[3].(*WithdrawFailedDueToInsufficientFunds)
	assert.True(t, ok)
	assert.Equal(t, withdrawMoney.ID, eventType.ID)
	assert.Equal(t, withdrawMoney.Amount, eventType.Amount)
	assert.Equal(t, 0, eventType.Balance)
}

func Test_CloseAccount(t *testing.T) {
	t.Parallel()

	// Given
	eventStore := Seacrest.NewEventStore()
	checkingAccountService := New(eventStore)
	id := "ABCD"
	name := "Alex Gemmell"
	accountWasOpened := AccountWasOpened{
		ID:   id,
		Name: name,
	}
	moneyWasDeposited := MoneyWasDeposited{
		ID:     id,
		Amount: 1099,
	}
	moneyWasWithdrawn := MoneyWasWithdrawn{
		ID:     id,
		Amount: 1099,
	}
	historicalEvents := append([]Event{}, accountWasOpened, moneyWasDeposited, moneyWasWithdrawn)
	err := checkingAccountService.PersistEvents(historicalEvents...)
	assert.Nil(t, err)

	// When
	closeAccount := CloseAccount{
		ID:   id,
	}
	err = checkingAccountService.HandleCommand(closeAccount)
	assert.Nil(t, err)

	// Then
	events, err := checkingAccountService.GetAllEvents()
	assert.Nil(t, err)
	assert.Len(t, events, 4)

	eventType, ok := events[3].(*AccountWasClosed)
	assert.True(t, ok)
	assert.Equal(t, closeAccount.ID, eventType.ID)
}
