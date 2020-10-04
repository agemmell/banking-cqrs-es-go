package main

import (
	"fmt"
	"github.com/agemmell/banking-cqrs-es-go/CheckingAccountService"
	"github.com/agemmell/banking-cqrs-es-go/Projections"
	"github.com/agemmell/banking-cqrs-es-go/Seacrest"
	"os"
	"time"
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

	fmt.Print("[generating events...")
	timer := time.Now()
	err := cas.AutoGenerateEvents()
	if err != nil {
		fmt.Printf("error %+v", err)
		os.Exit(1)
	}
	diff := time.Now().Sub(timer)
	fmt.Printf(" done] (%s)\n", diff.String())

	fmt.Print("[writing events to file...")
	timer = time.Now()
	err = cas.WriteEventsToFile("generated_events.txt")
	if err != nil {
		handleErrorAndExit(err)
	}
	diff = time.Now().Sub(timer)
	fmt.Printf(" done] (%s)\n", diff.String())
}

func LoadCheckingAccountEvents(eventStore *Seacrest.EventStore) {
	err := eventStore.LoadEventsFromFile("generated_events.txt")
	if err != nil {
		handleErrorAndExit(err)
	}
}

func RunProjections() {
	eventStore := Seacrest.NewEventStore()
	fmt.Print("[loading events into event store...")
	timer := time.Now()
	LoadCheckingAccountEvents(eventStore)
	diff := time.Now().Sub(timer)
	fmt.Printf(" done] (%s)\n", diff.String())

	// Run a projection
	fmt.Printf("\nRunning projections:\n")

	timer = time.Now()
	err := Projections.TotalBankFunds(eventStore)
	if err != nil {
		handleErrorAndExit(err)
	}
	diff = time.Now().Sub(timer)
	fmt.Printf("[TotalBankFunds done] (%s)\n\n", diff.String())

	timer = time.Now()
	err = Projections.OpenClosedAccounts(eventStore)
	if err != nil {
		handleErrorAndExit(err)
	}
	diff = time.Now().Sub(timer)
	fmt.Printf("[OpenClosedAccounts done] (%s)\n\n", diff.String())

	timer = time.Now()
	err = Projections.HighestBalanceOwners(eventStore)
	if err != nil {
		handleErrorAndExit(err)
	}
	diff = time.Now().Sub(timer)
	fmt.Printf("[HighestBalanceOwners done] (%s)\n\n", diff.String())

	timer = time.Now()
	err = Projections.TotalBalancePerMonth(eventStore)
	if err != nil {
		handleErrorAndExit(err)
	}
	diff = time.Now().Sub(timer)
	fmt.Printf("[TotalBalancePerMonth done] (%s)\n\n", diff.String())
}

func handleErrorAndExit(err error) {
	fmt.Printf("error %+v", err)
	os.Exit(1)
}
