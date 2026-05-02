package service

import (
	"errors"
	"slices"
	"virtual-wallet/internal/models/transaction"
	"virtual-wallet/internal/models/wallet"
	"virtual-wallet/internal/repository"
)

type WalletRepository interface {
	GetWalletsByProfileID(profileID int64) ([]*wallet.Wallet, error)
	CreateWallet(profileID int64, currency string) (int64, error)
	AddFunds(idempotencyKey string, walletID int64, profileID int64, amount int64) error
	TransferFunds(idempotencyKey string, profileID int64, fromWalletID int64, toWalletID int64, amount int64) error
	GetTransactionsHistory(walletID int64) ([]*transaction.Transaction, error)
}

type WalletService struct {
	repository WalletRepository
}

func NewWalletService(repository WalletRepository) *WalletService {
	return &WalletService{repository: repository}
}

var ErrInvalidAmount = errors.New("amount is less or equal than 0")
var ErrSameWallet = errors.New("fromWalletID and toWalletID is the same")
var ErrUnauthorizedAccess = errors.New("wallet does not belong to user")

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

func (r *WalletService) AddFunds(idempotencyKey string, walletID int64, profileID int64, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	errAddFunds := r.repository.AddFunds(idempotencyKey, walletID, profileID, amount)

	if errAddFunds != nil {
		if errors.Is(errAddFunds, repository.ErrIdempotentRequest) {
			return nil
		}
		
		return errAddFunds
	}

	return nil
}

func (r *WalletService) TransferFunds(idempotencyKey string, profileID int64, fromWalletID int64, toWalletID int64,
	amount int64) error {
	if fromWalletID == toWalletID {
		return ErrSameWallet
	}

	if amount <= 0 {
		return ErrInvalidAmount
	}

	errTransferFunds := r.repository.TransferFunds(idempotencyKey, profileID, fromWalletID, toWalletID, amount)

	if errTransferFunds != nil {
		if errors.Is(errTransferFunds, repository.ErrIdempotentRequest) {
			return nil
		}

		return errTransferFunds
	}

	return nil
}

func (r *WalletService) GetTransactionsHistory(profileID int64, walletID int64) ([]*transaction.Transaction, error) {
	wallets, errWallets := r.GetWalletsByProfileID(profileID)

	if errWallets != nil {
		return nil, errWallets
	}

	var walletIn bool

	for i := 0; i < len(wallets); i++ {
		if wallets[i].ID == walletID {
			walletIn = true
		}
	}

	if !walletIn {
		return nil, ErrUnauthorizedAccess
	}

	history, errHistory := r.repository.GetTransactionsHistory(walletID)

	if errHistory != nil {
		return nil, errHistory
	}

	return history, nil
}
