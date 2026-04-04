package service

import (
	"virtual-wallet/internal/models/user"
	"virtual-wallet/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (r *UserService) RegisterNewUser(firstName string, lastName string, email string, login string, password string) (int64, error) {
	passwordHash, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if errHash != nil {
		return 0, errHash
	}

	profile := &user.UserProfile{FirstName: firstName, LastName: lastName, Email: email}
	creds := &user.UserCredentials{Username: login, PasswordHash: string(passwordHash)}

	registerUserID, errRegister := r.repository.RegisterUser(profile, creds)
	if errRegister != nil {
		return 0, errRegister
	}

	return registerUserID, nil
}
