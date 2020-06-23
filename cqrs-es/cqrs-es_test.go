package cqrs_es

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CQRSESService_CreateMessageOfType(t *testing.T) {
	t.Parallel()

	testMessageType := "TestMessageType"

	s := &CQRSESService{}
	got, err := s.CreateMessageOfType(testMessageType)
	assert.Nil(t, err)
	assert.Equal(t, testMessageType, got.MessageType())
}

func Test_CQRSESService_CreateUUID(t *testing.T) {
	t.Parallel()

	CQRSESService := CQRSESService{}
	got, err := CQRSESService.CreateUUID()
	assert.Nil(t, err)
	assert.Regexp(t, `(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[0-9A-F]{4}-[0-9A-F]{12}$`, got.String())
}

func Test_CQRSESService_CreateUUIDString(t *testing.T) {
	t.Parallel()

	CQRSESService := CQRSESService{}
	got, err := CQRSESService.CreateUUIDString()
	assert.Nil(t, err)
	assert.True(t, len(got) > 0)
}

func Test_Message_MessageID(t *testing.T) {
	t.Parallel()

	testMessage := Message{"test-id", "test-type"}
	got := testMessage.MessageID()
	assert.Equal(t, "test-id", got)
}

func Test_Message_MessageType(t *testing.T) {
	t.Parallel()

	testMessage := Message{"test-id", "test-type"}
	got := testMessage.MessageType()
	assert.Equal(t, "test-type", got)
}
