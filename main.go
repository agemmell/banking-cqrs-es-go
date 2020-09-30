package main

import (
	"fmt"
	"github.com/agemmell/banking-cqrs-es-go/CheckingAccountService"
	"github.com/agemmell/banking-cqrs-es-go/Projections"
	"github.com/agemmell/banking-cqrs-es-go/Seacrest"
	"os"
)

func main() {
	// Generate some Checking Account events
	GenerateCheckingAccountEvents()

	// Run the projections
	RunProjections()
}

func GenerateCheckingAccountEvents() {
	eventStore := Seacrest.NewEventStore()
	cas := CheckingAccountService.New(eventStore)
	err := cas.GenerateCheckingAccountEvents()
	if err != nil {
		fmt.Printf("error %+v", err)
		os.Exit(1)
	}
	err = cas.WriteEventsToFile("generated_events.txt")
	if err != nil {
		handleErrorAndExit(err)
	}
}

func LoadCheckingAccountEvents(eventStore *Seacrest.EventStore) {
	err := eventStore.LoadEventsFromFile("generated_events.txt")
	if err != nil {
		handleErrorAndExit(err)
	}
}

func RunProjections() {
	eventStore := Seacrest.NewEventStore()
	LoadCheckingAccountEvents(eventStore)

	// Run a projection
	err := Projections.TotalBankFunds(eventStore)
	if err != nil {
		handleErrorAndExit(err)
	}

	err = Projections.OpenClosedAccounts(eventStore)
	if err != nil {
		handleErrorAndExit(err)
	}

	err = Projections.HighestBalanceOwners(eventStore)
	if err != nil {
		handleErrorAndExit(err)
	}
}

func handleErrorAndExit(err error) {
	fmt.Printf("error %+v", err)
	os.Exit(1)
}
