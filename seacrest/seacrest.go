package seacrest

import uuid "github.com/nu7hatch/gouuid"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . MessageDescriber
type MessageDescriber interface {
	MessageID() string
	MessageType() string
}

//counterfeiter:generate . CommandHandler
type CommandHandler interface {
	HandleCommand() error
}

//counterfeiter:generate . Seacrest
type Seacrest interface {
	CreateUUID() (*uuid.UUID, error)
	CreateUUIDString() (string, error)
	CreateMessageOfType(messageType string) (MessageDescriber, error)
}

type Message struct {
	messageID   string
	messageType string
}

type Service struct{}

func (m Message) MessageID() string {
	return m.messageID
}

func (m Message) MessageType() string {
	return m.messageType
}

func (s *Service) CreateUUID() (*uuid.UUID, error) {
	uuid4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return uuid4, nil
}

func (s *Service) CreateUUIDString() (string, error) {
	uuid4, err := s.CreateUUID()
	if err != nil {
		return "", err
	}

	return uuid4.String(), nil
}

func (s *Service) CreateMessageOfType(messageType string) (MessageDescriber, error) {
	messageID, err := s.CreateUUIDString()
	if err != nil {
		return nil, err
	}

	return Message{
		messageID,
		messageType,
	}, nil
}