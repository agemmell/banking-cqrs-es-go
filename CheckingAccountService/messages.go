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

func (c OpenAccount) isCommand()  {}
func (c DepositMoney) isCommand() {}

// Events

type AccountWasOpened struct {
	ID      string
	Name    string
	version uint
}
type MoneyWasDeposited struct {
	ID      string
	Amount  int
	version uint
}
type MoneyWasWithdrawn struct {
	ID      string
	Amount  int
	version uint
}
type WithdrawFailedDueToInsufficientFunds struct {
	ID      string
	Amount  int
	version uint
}
type AccountWasClosed struct {
	ID      string
	version uint
}

func (e AccountWasOpened) isEvent()                     {}
func (e MoneyWasDeposited) isEvent()                    {}
func (e MoneyWasWithdrawn) isEvent()                    {}
func (e WithdrawFailedDueToInsufficientFunds) isEvent() {}
func (e AccountWasClosed) isEvent()                     {}

// AggregateID() and Version() satisfy the Seacrest.Event interface. This is needed to cast between the two when
// persisting/getting events from the Seacrest event store
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

func (e AccountWasOpened) Version() uint {
	return e.version // TODO I need to add the version to the event on creation so it gets stored in the event store - then tests wil lpass?
}
func (e MoneyWasDeposited) Version() uint {
	return e.version
}
func (e MoneyWasWithdrawn) Version() uint {
	return e.version
}
func (e WithdrawFailedDueToInsufficientFunds) Version() uint {
	return e.version
}
func (e AccountWasClosed) Version() uint {
	return e.version
}
