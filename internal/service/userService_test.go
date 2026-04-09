package service

import (
	"testing"
	"virtual-wallet/internal/models/user"
)

type MockUserRepository struct{}

func (m *MockUserRepository) RegisterUser(profile *user.UserProfile, creds *user.UserCredentials) (int64, error) {
	return 1, nil
}

func TestRegisterNewUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	svc := NewUserService(mockRepo)
	registerUserID, errRegister := svc.RegisterNewUser("John", "Doe", "john@doe.com", "johndoe", "johndoepassword")
	if errRegister != nil {
		t.Errorf("errRegister: %s", errRegister)
	}

	if registerUserID != 1 {
		t.Errorf("registerUserID is not 1")
	}

}
