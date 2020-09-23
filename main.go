package main

import (
	"fmt"
	"github.com/agemmell/banking-cqrs-es-go/CheckingAccountService"
	"github.com/agemmell/banking-cqrs-es-go/Seacrest"
	"os"
)

func main() {
	eventStore := Seacrest.NewEventStore()
	cas := CheckingAccountService.New(eventStore)
	err := cas.GenerateCheckingAccountEvents()
	if err != nil {
		fmt.Printf("error %+v", err)
		os.Exit(1)
	}
}
