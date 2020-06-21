package cqrs_es

import uuid "github.com/nu7hatch/gouuid"

type MessageDescriber interface {
	MessageID() string
	MessageType() string
}

type CommandHandler interface {
	HandleCommand() error
}

type CQRSES interface {
	CreateMessageIDUUIDv4() (string, error)
	CreateMessageOfType(messageType string) (MessageDescriber, error)
}

type Message struct {
	messageID   string
	messageType string
}

type CQRSESService struct{}

func (m *Message) MessageID() string {
	return m.messageID
}

func (m *Message) MessageType() string {
	return m.messageType
}

func (s *CQRSESService) CreateMessageIDUUIDv4() (string, error) {
	uuid4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return uuid4.String(), nil
}

func (s *CQRSESService) CreateMessageOfType(messageType string) (MessageDescriber, error) {
	messageID, err := s.CreateMessageIDUUIDv4()
	if err != nil {
		return nil, err
	}

	return &Message{
		messageID,
		messageType,
	}, nil
}
