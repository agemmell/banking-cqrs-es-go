package Projections

import (
	"encoding/json"
	"fmt"
	"github.com/agemmell/banking-cqrs-es-go/CheckingAccountService"
	"github.com/agemmell/banking-cqrs-es-go/Seacrest"
	"sort"
)

// Total bank funds (sum of all account balances)
func TotalBankFunds(eventStore *Seacrest.EventStore) error {

	totalBankFunds := 0

	for _, envelope := range eventStore.GetAllEvents() {
		switch envelope.EventType {
		case "MoneyWasDeposited":
			moneyWasDeposited := CheckingAccountService.MoneyWasDeposited{}
			err := json.Unmarshal(envelope.Payload, &moneyWasDeposited)
			if err != nil {
				return err
			}
			totalBankFunds += moneyWasDeposited.Amount
		case "MoneyWasWithdrawn":
			moneyWasWithdrawn := CheckingAccountService.MoneyWasWithdrawn{}
			err := json.Unmarshal(envelope.Payload, &moneyWasWithdrawn)
			if err != nil {
				return err
			}
			totalBankFunds -= moneyWasWithdrawn.Amount
		}
	}

	fmt.Printf("Total Banks Funds = $%d\n", totalBankFunds)

	return nil
}

// Number of open & closed accounts
func OpenClosedAccounts(eventStore *Seacrest.EventStore) error {

	openAccounts, closedAccounts := 0, 0

	for _, envelope := range eventStore.GetAllEvents() {
		switch envelope.EventType {
		case "AccountWasOpened":
			accountWasOpened := CheckingAccountService.AccountWasOpened{}
			err := json.Unmarshal(envelope.Payload, &accountWasOpened)
			if err != nil {
				return err
			}
			openAccounts += 1
		case "AccountWasClosed":
			accountWasClosed := CheckingAccountService.AccountWasClosed{}
			err := json.Unmarshal(envelope.Payload, &accountWasClosed)
			if err != nil {
				return err
			}
			openAccounts -= 1
			closedAccounts += 1
		}
	}

	fmt.Printf("Open Accounts = %d\n", openAccounts)
	fmt.Printf("Closed Accounts = %d\n", closedAccounts)

	return nil
}

type AccountBalance struct {
	ID      string
	Name    string
	Balance int
}

// Top 10 "highest balance" account owners
func HighestBalanceOwners(eventStore *Seacrest.EventStore) error {

	var accountBalances = map[string]AccountBalance{}
	var topTen []AccountBalance

	for _, envelope := range eventStore.GetAllEvents() {
		switch envelope.EventType {
		case "AccountWasOpened":
			accountWasOpened := CheckingAccountService.AccountWasOpened{}
			err := json.Unmarshal(envelope.Payload, &accountWasOpened)
			if err != nil {
				return err
			}
			accountBalance := AccountBalance{
				ID:      accountWasOpened.ID,
				Name:    accountWasOpened.Name,
				Balance: 0,
			}
			accountBalances[accountWasOpened.ID] = accountBalance
			topTen = sortTopTen(topTen, accountBalance)

		case "MoneyWasDeposited":
			moneyWasDeposited := CheckingAccountService.MoneyWasDeposited{}
			err := json.Unmarshal(envelope.Payload, &moneyWasDeposited)
			if err != nil {
				return err
			}
			accountBalance := accountBalances[moneyWasDeposited.ID]
			accountBalance.Balance += moneyWasDeposited.Amount
			accountBalances[moneyWasDeposited.ID] = accountBalance
			topTen = sortTopTen(topTen, accountBalance)

		case "MoneyWasWithdrawn":
			moneyWasWithdrawn := CheckingAccountService.MoneyWasWithdrawn{}
			err := json.Unmarshal(envelope.Payload, &moneyWasWithdrawn)
			if err != nil {
				return err
			}
			accountBalance := accountBalances[moneyWasWithdrawn.ID]
			accountBalance.Balance = moneyWasWithdrawn.Balance
			accountBalances[moneyWasWithdrawn.ID] = accountBalance
			topTen = sortTopTen(topTen, accountBalance)
		}
	}

	fmt.Println("Top Ten Balances:")
	for i, account := range topTen {
		fmt.Printf("%d. %d (%s: %s)\n", i+1, account.Balance, account.Name, account.ID)
	}

	return nil
}

func sortTopTen(topTen []AccountBalance, accountBalance AccountBalance) []AccountBalance {
	alreadyExists := false
	for i, account := range topTen {
		if account.ID == accountBalance.ID {
			topTen[i] = accountBalance
			alreadyExists = true
			break
		}
	}

	if !alreadyExists {
		topTen = append(topTen, accountBalance)
	}

	sort.Slice(topTen, func(i, j int) bool {
		return topTen[i].Balance > topTen[j].Balance
	})

	l := len(topTen)
	if l > 10 {
		l = 10
	}
	return topTen[0:l]
}

