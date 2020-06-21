package banking

import cqrses "github.com/agemmell/banking-cqrs-es-go/cqrs-es"

const OpenAccountMessageType = "OpenAccount"
const AccountWasOpenedMessageType = "AccountWasOpened"

type OpenAccount struct {
	cqrses.MessageDescriber
	accountID string
	name      string
}

func (c *OpenAccount) AccountID() string {
	return c.accountID
}

func (c *OpenAccount) Name() string {
	return c.name
}

type AccountWasOpened struct {
	cqrses.MessageDescriber
	accountID string
	name      string
}

func (e AccountWasOpened) AccountID() string {
	return e.accountID
}

func (e AccountWasOpened) Name() string {
	return e.name
}
