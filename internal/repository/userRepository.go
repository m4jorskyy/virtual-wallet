package repository

import (
	"database/sql"
	"virtual-wallet/internal/models/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) RegisterUser(profile *user.UserProfile, creds *user.UserCredentials) (int64, error) {
	tx, errTx := r.db.Begin()
	if errTx != nil {
		return 0, errTx
	}

	var returnedID int64

	errReturnID := tx.QueryRow("INSERT INTO user_profiles (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING user_profiles.id", profile.FirstName, profile.LastName, profile.Email).Scan(&returnedID)

	if errReturnID != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return 0, errRollback
		}
		return 0, errReturnID
	}

	_, errExec := tx.Exec("INSERT INTO user_credentials (username, password_hash, profile_id) VALUES ($1, $2, $3)", creds.Username, creds.PasswordHash, returnedID)

	if errExec != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return 0, errRollback
		}
		return 0, errExec
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return 0, errCommit
	}

	return returnedID, nil
}
