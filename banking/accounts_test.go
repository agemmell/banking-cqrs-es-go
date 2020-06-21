package banking

import (
	"testing"

	uuid "github.com/nu7hatch/gouuid"

	cqrsesfakes "github.com/agemmell/banking-cqrs-es-go/cqrs-es/cqrs-esfakes"
)

func Test_CreateOpenAccount(t *testing.T) {
	t.Parallel()

	// Given
	fakeCQRSESService := cqrsesfakes.FakeCQRSES{}

	fakeUUID := "fake-uuid-string"
	fakeMessage := cqrsesfakes.FakeMessageDescriber{}
	fakeMessage.MessageIDReturns(fakeUUID)
	fakeMessage.MessageTypeReturns(OpenAccountMessageType)
	fakeCQRSESService.CreateMessageOfTypeReturns(&fakeMessage, nil)

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
	if got.MessageID() != fakeUUID {
		t.Errorf("Got %v, want %v", got.MessageID(), fakeUUID)
	}
	if got.MessageType() != OpenAccountMessageType {
		t.Errorf("Got %v, want %v", got.MessageType(), OpenAccountMessageType)
	}
}

func Test_CreateOpenAccountWithUUID(t *testing.T) {
	t.Parallel()

	// Given
	fakeCQRSESService := cqrsesfakes.FakeCQRSES{}

	fakeUUID := "fake-uuid-string"
	fakeMessage := cqrsesfakes.FakeMessageDescriber{}
	fakeMessage.MessageIDReturns(fakeUUID)
	fakeMessage.MessageTypeReturns(OpenAccountMessageType)
	fakeCQRSESService.CreateMessageOfTypeReturns(&fakeMessage, nil)

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
	if got.MessageID() != fakeUUID {
		t.Errorf("Got %v, want %v", got.MessageID(), fakeUUID)
	}
	if got.MessageType() != OpenAccountMessageType {
		t.Errorf("Got %v, want %v", got.MessageType(), OpenAccountMessageType)
	}
}

func Test_CreateAccountWasOpened(t *testing.T) {
	t.Parallel()

	// Given
	fakeCQRSESService := cqrsesfakes.FakeCQRSES{}

	fakeUUID := "fake-uuid-string"
	fakeMessage := cqrsesfakes.FakeMessageDescriber{}
	fakeMessage.MessageIDReturns(fakeUUID)
	fakeMessage.MessageTypeReturns(AccountWasOpenedMessageType)
	fakeCQRSESService.CreateMessageOfTypeReturns(&fakeMessage, nil)

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
	if got.MessageID() != fakeUUID {
		t.Errorf("Got %v, want %v", got.MessageID(), fakeUUID)
	}
	if got.MessageType() != AccountWasOpenedMessageType {
		t.Errorf("Got %v, want %v", got.MessageType(), AccountWasOpenedMessageType)
	}
}
