package repository

import (
	"database/sql"
	"virtual-wallet/internal/models/wallet"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

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
