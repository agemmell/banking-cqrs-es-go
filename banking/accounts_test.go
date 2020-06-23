package banking

import (
	"testing"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"

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
	assert.Nil(t, err)

	// Then
	assert.Equal(t, name, got.Name())
	assert.Equal(t, accountID, got.AccountID())
	assert.Equal(t, fakeUUID, got.MessageID())
	assert.Equal(t, OpenAccountMessageType, got.MessageType())
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
	assert.Nil(t, err)

	// Then
	assert.Equal(t, name, got.Name())
	assert.Equal(t, accountID.String(), got.AccountID())
	assert.Equal(t, fakeUUID, got.MessageID())
	assert.Equal(t, OpenAccountMessageType, got.MessageType())
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
	assert.Nil(t, err)

	// Then
	assert.Equal(t, name, got.Name())
	assert.Equal(t, accountID, got.AccountID())
	assert.Equal(t, fakeUUID, got.MessageID())
	assert.Equal(t, AccountWasOpenedMessageType, got.MessageType())
}
