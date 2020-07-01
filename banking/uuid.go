package banking

import uuid "github.com/nu7hatch/gouuid"

type UuidService struct{}

func (uuidService *UuidService) CreateUUIDString() (string, error) {
	uuidV4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return uuidV4.String(), nil
}
