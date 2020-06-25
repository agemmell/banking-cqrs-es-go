package banking

import (
	"errors"

	"github.com/agemmell/banking-cqrs-es-go/seacrest"
)

type Banking struct {
	eventStore     seacrest.StoresEvents
}

func NewService(eventStore cqrses.StoresEvents) Banking {
	return Banking{eventStore}
}

func (b *Banking) HandleCommand(c cqrses.MessageDescriber) {
	// todo
}
