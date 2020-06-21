package cqrs_es

import "testing"

func Test_CreateMessageIDUUIDv4(t *testing.T) {
	CQRSESService := CQRSESService{}
	got, err := CQRSESService.CreateMessageIDUUIDv4()
	if err != nil {
		t.Errorf("CreateMessageIDUUIDv4() error = %v", err)
	}

	if len(got) < 1 {
		t.Errorf("CreateMessageIDUUIDv4() returned empty string: %v", got)
	}
}