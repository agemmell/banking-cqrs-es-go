package banking

const OpenAccountMessageType = "OpenAccount"
const AccountWasOpenedMessageType = "AccountWasOpened"

type Message struct {
	messageID   string
	messageType string
}

func (m *Message) MessageID() string {
	return m.messageID
}

func (m *Message) MessageType() string {
	return m.messageType
}

// Commands

type OpenAccount struct {
	Message
	accountID string
	name      string
}

func (c *OpenAccount) AccountID() string {
	return c.accountID
}

func (c *OpenAccount) Name() string {
	return c.name
}

// Events

type AccountWasOpened struct {
	Message
	accountID string
	name      string
}

func (e AccountWasOpened) AccountID() string {
	return e.accountID
}

func (e AccountWasOpened) Name() string {
	return e.name
}
