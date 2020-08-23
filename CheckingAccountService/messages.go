package CheckingAccountService

// Commands

type OpenAccount struct {
	AccountID string
	Name      string
}

type DepositMoney struct {
	AccountID     string
	DepositAmount int
}

func (c OpenAccount) isCommand()  {}
func (c DepositMoney) isCommand() {}

// Events

type AccountWasOpened struct {
	AccountID string
	Name      string
}

type MoneyWasDeposited struct {
	AccountID       string
	DepositedAmount int
}

func (e AccountWasOpened) isEvent()  {}
func (e MoneyWasDeposited) isEvent() {}
