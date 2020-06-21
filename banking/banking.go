package banking

import cqrses "github.com/agemmell/banking-cqrs-es-go/cqrs-es"

type Banking struct {
	eventStore cqrses.StoresEvents
}

func NewService(eventStore cqrses.StoresEvents) Banking {
	return Banking{eventStore}
}

func (b *Banking) HandleCommand(c cqrses.MessageDescriber) {

}
