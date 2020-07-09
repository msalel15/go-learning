package user

import (
	"testing"
)

type mockUserService struct {
}

func (m *mockUserService) CreateUser(user UserDTO) error {
	panic("implement me")
}

func (m *mockUserService) IsExistUsernameAndPassword(username, password string) (bool, error) {
	if username == "exist" {
		return true, nil
	}
	return false, nil
}

func (m mockUserService) IsExistByUsername(username string) (bool, error) {
	panic("implement me")
}

func TestSignUpHandler(t *testing.T) {
	user := User{server: nil, service: new(mockUserService)}

	user.signUpHandler(nil)
}
