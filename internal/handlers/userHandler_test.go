package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"virtual-wallet/internal/models/user"
	"virtual-wallet/internal/service"
)

type MockUserRepository struct{}

func (m *MockUserRepository) RegisterUser(profile *user.UserProfile, creds *user.UserCredentials) (int64, error) {
	return 1, nil
}

func TestUserHandler_RegisterUser(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/register/", strings.NewReader("{\n  \"first_name\": \"John\",\n  \"last_name\": \"Doe\",\n  \"email\": \"john.doe@gmail.com\",\n  \"username\": \"johndoe\",\n  \"password\": \"johndoe123\"\n}"))

	recorder := httptest.NewRecorder()

	mockRepo := &MockUserRepository{}
	svc := service.NewUserService(mockRepo)
	handler := NewUserHandler(svc)
	handler.RegisterUser(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Register Code is not 201")
	}
}
