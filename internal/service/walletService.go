package service

import (
	"errors"
	"slices"
	"virtual-wallet/internal/models/wallet"
)

type WalletRepository interface {
	GetWalletsByProfileID(profileID int64) ([]*wallet.Wallet, error)
	CreateWallet(profileID int64, currency string) (int64, error)
	AddFunds(walletID int64, profileID int64, amount int64) error
}

type WalletService struct {
	repository WalletRepository
}

func NewWalletService(repository WalletRepository) *WalletService {
	return &WalletService{repository: repository}
}

var ErrInvalidAmount = errors.New("amount is less or equal than 0")

func (r *WalletService) GetWalletsByProfileID(profileID int64) ([]*wallet.Wallet, error) {
	wallets, errWallets := r.repository.GetWalletsByProfileID(profileID)

	if errWallets != nil {
		return nil, errWallets
	}

	return wallets, nil
}

func (r *WalletService) CreateWallet(profileID int64, currency string) (int64, error) {
	availableCurrencies := []string{"PLN", "EUR", "USD", "GBP"}

	if !slices.Contains(availableCurrencies, currency) {
		return 0, errors.New("unavailable currency")
	}

	returnedWalletID, errCreateWallet := r.repository.CreateWallet(profileID, currency)

	if errCreateWallet != nil {
		return 0, errCreateWallet
	}

	return returnedWalletID, nil
}

func (r *WalletService) AddFunds(walletID int64, profileID int64, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	errAddFunds := r.repository.AddFunds(walletID, profileID, amount)

	if errAddFunds != nil {
		return errAddFunds
	}

	return nil
}
