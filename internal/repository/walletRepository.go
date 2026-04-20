package repository

import (
	"database/sql"
	"errors"
	"time"
	"virtual-wallet/internal/models/wallet"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

var ZeroRowsAffectedError error = errors.New("0 rows affected")

func (r *WalletRepository) GetWalletsByProfileID(profileID int64) ([]*wallet.Wallet, error) {
	var wallets []*wallet.Wallet

	rows, errWallets := r.db.Query("SELECT id, profile_id, balance, currency FROM wallet WHERE profile_id = $1", profileID)

	if errWallets != nil {
		return nil, errWallets
	}

	defer func(rows *sql.Rows) {
		errCloseRows := rows.Close()
		if errCloseRows != nil {
			return
		}
	}(rows)

	for rows.Next() {
		w := &wallet.Wallet{}

		errScan := rows.Scan(&w.ID, &w.ProfileID, &w.Balance, &w.Currency)

		if errScan != nil {
			return nil, errScan
		}

		wallets = append(wallets, w)
	}

	if errRows := rows.Err(); errRows != nil {
		return nil, errRows
	}

	return wallets, nil
}

func (r *WalletRepository) CreateWallet(profileID int64, currency string) (int64, error) {
	var returnedWalletID int64

	errCreateWallet := r.db.QueryRow("INSERT INTO wallet (profile_id, currency) VALUES ($1, $2) RETURNING id", profileID, currency).Scan(&returnedWalletID)
	if errCreateWallet != nil {
		return 0, errCreateWallet
	}

	return returnedWalletID, nil
}

func (r *WalletRepository) AddFunds(walletID int64, profileID int64, amount int64) error {
	tx, errTx := r.db.Begin()
	if errTx != nil {
		return errTx
	}

	defer func(tx *sql.Tx) {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return
		}
	}(tx)

	rows, errAddFunds := tx.Exec("UPDATE wallet SET balance = balance + $1 WHERE id = $2 AND profile_id = $3", amount, walletID, profileID)

	if errAddFunds != nil {
		return errAddFunds
	}

	affected, errRowsAffected := rows.RowsAffected()

	if errRowsAffected != nil {
		return errRowsAffected
	}

	if affected == 0 {
		return ZeroRowsAffectedError
	}

	_, errReturnedTransactionID := tx.Exec("INSERT INTO transactions (from_wallet_id, to_wallet_id, amount, created_at, type) VALUES ($1, $2, $3, $4, $5)", 0, walletID, amount, time.Now(), "DEPOSIT")

	if errReturnedTransactionID != nil {
		return errReturnedTransactionID
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (r *WalletRepository) TransferFunds(profileID int64, fromWalletID int64, toWalletID int64, amount int64) error {
	tx, errTx := r.db.Begin()
	if errTx != nil {
		return errTx
	}

	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	rows, errSubtractFunds := tx.Exec("UPDATE wallet SET balance = balance - $1 WHERE id = $2 AND profile_id = $3 AND balance - $1 >= 0", amount, fromWalletID, profileID)

	if errSubtractFunds != nil {
		return errSubtractFunds
	}

	affected, errRowsAffected := rows.RowsAffected()

	if errRowsAffected != nil {
		return errRowsAffected
	}

	if affected == 0 {
		return ZeroRowsAffectedError
	}

	_, errAddFunds := tx.Exec("UPDATE wallet SET balance = balance + $1 WHERE id = $2", amount, toWalletID)

	if errAddFunds != nil {
		return errAddFunds
	}

	_, errReturnedTransactionID := tx.Exec("INSERT INTO transactions (from_wallet_id, to_wallet_id, amount, created_at, type) VALUES ($1, $2, $3, $4, $5)", fromWalletID, toWalletID, amount, time.Now(), "TRANSFER")

	if errReturnedTransactionID != nil {
		return errReturnedTransactionID
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}
