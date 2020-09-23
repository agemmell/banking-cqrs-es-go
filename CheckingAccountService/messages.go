package CheckingAccountService

// Commands

type OpenAccount struct {
	ID   string
	Name string
}
type DepositMoney struct {
	ID     string
	Amount int
}
type WithdrawMoney struct {
	ID     string
	Amount int
}
type CloseAccount struct {
	ID string
}

func (c OpenAccount) isCommand()   {}
func (c DepositMoney) isCommand()  {}
func (c WithdrawMoney) isCommand() {}
func (c CloseAccount) isCommand()  {}

// Events

type AccountWasOpened struct {
	ID   string
	Name string
}
type MoneyWasDeposited struct {
	ID     string
	Amount int
}
type MoneyWasWithdrawn struct {
	ID      string
	Amount  int
	Balance int
}
type WithdrawFailedDueToInsufficientFunds struct {
	ID      string
	Amount  int
	Balance int
}
type AccountWasClosed struct {
	ID string
}

func (e AccountWasOpened) AggregateID() string {
	return e.ID
}
func (e MoneyWasDeposited) AggregateID() string {
	return e.ID
}
func (e MoneyWasWithdrawn) AggregateID() string {
	return e.ID
}
func (e WithdrawFailedDueToInsufficientFunds) AggregateID() string {
	return e.ID
}
func (e AccountWasClosed) AggregateID() string {
	return e.ID
}

const TypeAccountWasOpened = "AccountWasOpened"
const TypeMoneyWasDeposited = "MoneyWasDeposited"
const TypeMoneyWasWithdrawn = "MoneyWasWithdrawn"
const TypeWithdrawFailedDueToInsufficientFunds = "WithdrawFailedDueToInsufficientFunds"
const TypeAccountWasClosed = "AccountWasClosed"

func (e AccountWasOpened) EventType() string {
	return TypeAccountWasOpened
}
func (e MoneyWasDeposited) EventType() string {
	return TypeMoneyWasDeposited
}
func (e MoneyWasWithdrawn) EventType() string {
	return TypeMoneyWasWithdrawn
}
func (e WithdrawFailedDueToInsufficientFunds) EventType() string {
	return TypeWithdrawFailedDueToInsufficientFunds
}
func (e AccountWasClosed) EventType() string {
	return TypeAccountWasClosed
}
