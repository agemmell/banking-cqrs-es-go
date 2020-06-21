package banking

import (
	"testing"

	uuid "github.com/nu7hatch/gouuid"

	cqrses "github.com/agemmell/banking-cqrs-es-go/cqrs-es"
)

type FakeCQRSESService struct {
	fakeUUID string
}

type FakeMessage struct {
	messageID string
	messageType string
}

func (s *FakeCQRSESService) CreateMessageIDUUIDv4() (string, error) {
	return s.fakeUUID, nil
}

func (s *FakeCQRSESService) CreateMessageOfType(messageType string) (cqrses.MessageDescriber, error) {
	return &FakeMessage{
		s.fakeUUID,
		messageType,
	}, nil
}

func (m *FakeMessage) MessageID() string {
	return m.messageID
}

func (m *FakeMessage) MessageType() string {
	return m.messageType
}

func Test_CreateOpenAccount(t *testing.T) {

	// Given
	fakeCQRSESService := FakeCQRSESService{}
	fakeCQRSESService.fakeUUID = "fake-uuid-string"
	accountService := AccountService{&fakeCQRSESService}
	accountID := "test-account-id-string"
	name := "Alex Gemmell"

	// When
	got, err := accountService.CreateOpenAccount(accountID, name)
	if err != nil {
		t.Errorf("Error %v", err)
	}

	// Then
	if got.Name() != name {
		t.Errorf("Got %v, want %v", got.name, name)
	}
	if got.AccountID() != accountID {
		t.Errorf("Got %v, want %v", got.accountID, accountID)
	}
	if got.MessageID() != fakeCQRSESService.fakeUUID {
		t.Errorf("Got %v, want %v", got.MessageID(), fakeCQRSESService.fakeUUID)
	}
	if got.MessageType() != OpenAccountMessageType {
		t.Errorf("Got %v, want %v", got.MessageType(), OpenAccountMessageType)
	}
}

func Test_CreateOpenAccountWithUUID(t *testing.T) {

	// Given
	fakeCQRSESService := FakeCQRSESService{}
	accountService := AccountService{&fakeCQRSESService}
	accountID, _ := uuid.NewV4()
	name := "Alex Gemmell"

	// When
	got, err := accountService.CreateOpenAccountWithUUID(accountID, name)
	if err != nil {
		t.Errorf("Error %v", err)
	}

	// Then
	if got.Name() != name {
		t.Errorf("Got %v, want %v", got.name, name)
	}
	if got.AccountID() != accountID.String() {
		t.Errorf("Got %v, want %v", got.accountID, accountID.String())
	}
	if got.MessageID() != fakeCQRSESService.fakeUUID {
		t.Errorf("Got %v, want %v", got.MessageID(), fakeCQRSESService.fakeUUID)
	}
	if got.MessageType() != OpenAccountMessageType {
		t.Errorf("Got %v, want %v", got.MessageType(), OpenAccountMessageType)
	}
}

func Test_CreateAccountWasOpened(t *testing.T) {
	// Given
	fakeCQRSESService := FakeCQRSESService{}
	accountService := AccountService{&fakeCQRSESService}
	accountID := "test-account-id-string"
	name := "Alex Gemmell"

	// When
	got, err := accountService.CreateAccountWasOpened(accountID, name)
	if err != nil {
		t.Errorf("Error %v", err)
	}

	// Then
	if got.Name() != name {
		t.Errorf("Got %v, want %v", got.name, name)
	}
	if got.AccountID() != accountID {
		t.Errorf("Got %v, want %v", got.accountID, accountID)
	}
	if got.MessageID() != fakeCQRSESService.fakeUUID {
		t.Errorf("Got %v, want %v", got.MessageID(), fakeCQRSESService.fakeUUID)
	}
	if got.MessageType() != AccountWasOpenedMessageType {
		t.Errorf("Got %v, want %v", got.MessageType(), AccountWasOpenedMessageType)
	}
}
