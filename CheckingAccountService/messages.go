package CheckingAccountService

// Commands

type OpenAccount struct {
	AccountID string
	Name      string
}

type DepositMoney struct {
	AccountID string
	Amount    int
}

func (c OpenAccount) isCommand()  {}
func (c DepositMoney) isCommand() {}

// Events

type AccountWasOpened struct {
	AccountID string
	Name      string
}

type MoneyWasDeposited struct {
	AccountID string
	Amount    int
}

type MoneyWasWithdrawn struct {
	AccountID string
	Amount    int
}
type WithdrawFailedDueToInsufficientFunds struct {
	AccountID string
	Amount    int
}

func (e AccountWasOpened) isEvent()                     {}
func (e MoneyWasDeposited) isEvent()                    {}
func (e MoneyWasWithdrawn) isEvent()                    {}
func (e WithdrawFailedDueToInsufficientFunds) isEvent() {}
