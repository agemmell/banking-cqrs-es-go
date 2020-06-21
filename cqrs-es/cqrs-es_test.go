package cqrs_es

import (
	"testing"
)

func Test_CQRSESService_CreateMessageOfType(t *testing.T) {
	t.Parallel()

	type args struct {
		messageType string
	}

	testMessageType := "TestMessageType"

	tests := []struct {
		name    string
		args    args
		want    MessageDescriber
		wantErr bool
	}{
		{
			"happy path",
			args{
				testMessageType,
			},
			&Message{"test-uuid", testMessageType},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CQRSESService{}
			got, _ := s.CreateMessageOfType(tt.args.messageType)
			if got.MessageType() != tt.want.MessageType() {
				t.Errorf("MessageType() got = %v, want %v", got.MessageType(), tt.want.MessageType())
			}
		})
	}
}

func Test_CQRSESService_CreateUUID(t *testing.T) {
	t.Parallel()

	CQRSESService := CQRSESService{}
	got, err := CQRSESService.CreateUUID()
	if err != nil {
		t.Errorf("CreateUUID() error = %v", err)
	}

	if len(got) < 1 {
		t.Errorf("CreateUUID() returned empty string: %v", got)
	}
}

func Test_CQRSESService_CreateUUIDString(t *testing.T) {
	t.Parallel()

	CQRSESService := CQRSESService{}
	got, err := CQRSESService.CreateUUIDString()
	if err != nil {
		t.Errorf("CreateUUIDString() error = %v", err)
	}

	if len(got) < 1 {
		t.Errorf("CreateUUIDString() returned empty string: %v", got)
	}
}

func Test_Message_MessageID(t *testing.T) {
	t.Parallel()

	testMessage := Message{"test-id", "test-type"}
	got := testMessage.MessageID()
	if got != "test-id" {
		t.Errorf("MessageID() got %v, want %v", got, "test-id")
	}
}

func Test_Message_MessageType(t *testing.T) {
	t.Parallel()

	testMessage := Message{"test-id", "test-type"}
	got := testMessage.MessageType()
	if got != "test-type" {
		t.Errorf("MessageType() got %v, want %v", got, "test-type")
	}
}
