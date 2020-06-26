package banking

import (
	"testing"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"

	"github.com/agemmell/banking-cqrs-es-go/seacrest/seacrestfakes"
)

func Test_CreateOpenAccount(t *testing.T) {
	t.Parallel()

	// Given
	fakeSeacrestService := seacrestfakes.FakeSeacrest{}

	fakeUUID := "fake-uuid-string"
	fakeMessage := seacrestfakes.FakeMessageDescriber{}
	fakeMessage.MessageIDReturns(fakeUUID)
	fakeMessage.MessageTypeReturns(OpenAccountMessageType)
	fakeSeacrestService.CreateMessageOfTypeReturns(&fakeMessage, nil)

	accountService := AccountService{&fakeSeacrestService}
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
	fakeSeacrestService := seacrestfakes.FakeSeacrest{}

	fakeUUID := "fake-uuid-string"
	fakeMessage := seacrestfakes.FakeMessageDescriber{}
	fakeMessage.MessageIDReturns(fakeUUID)
	fakeMessage.MessageTypeReturns(OpenAccountMessageType)
	fakeSeacrestService.CreateMessageOfTypeReturns(&fakeMessage, nil)

	accountService := AccountService{&fakeSeacrestService}
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
	fakeSeacrestService := seacrestfakes.FakeSeacrest{}

	fakeUUID := "fake-uuid-string"
	fakeMessage := seacrestfakes.FakeMessageDescriber{}
	fakeMessage.MessageIDReturns(fakeUUID)
	fakeMessage.MessageTypeReturns(AccountWasOpenedMessageType)
	fakeSeacrestService.CreateMessageOfTypeReturns(&fakeMessage, nil)

	accountService := AccountService{&fakeSeacrestService}
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

func Test_Account_ApplyEventUnknownType(t *testing.T) {
	t.Parallel()

	// Given
	account := Account{}

	fakeMessageUnknownType := seacrestfakes.FakeMessageDescriber{}
	fakeMessageUnknownType.MessageTypeReturns("test-type")

	// When
	err := account.ApplyEvent(&fakeMessageUnknownType)

	// Then
	assert.Equal(t, "unknown event type test-type", err.Error())
}