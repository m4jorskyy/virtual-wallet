package service

import (
	"errors"
	"testing"
	"time"
	"virtual-wallet/internal/models/transaction"
	"virtual-wallet/internal/models/wallet"
	"virtual-wallet/internal/repository"
)

type MockWalletRepository struct{}

func (m *MockWalletRepository) GetWalletsByProfileID(profileID int64) ([]*wallet.Wallet, error) {
	var wallets []*wallet.Wallet
	w := &wallet.Wallet{}

	w.ID = 1
	w.ProfileID = profileID
	w.Balance = 0
	w.Currency = "PLN"

	wallets = append(wallets, w)

	return wallets, nil
}

func (m *MockWalletRepository) CreateWallet(profileID int64, currency string) (int64, error) {
	return 1, nil
}

func (m *MockWalletRepository) AddFunds(idempotencyKey string, walletID int64, profileID int64, amount int64) error {
	if idempotencyKey == "asdfghjkl1234567890" {
		return repository.ErrIdempotentRequest
	}

	return nil
}

func (m *MockWalletRepository) TransferFunds(idempotencyKey string, profileID int64, fromWalletID int64,
	toWalletID int64, amount int64) error {
	if idempotencyKey == "asdfghjkl1234567890" {
		return repository.ErrIdempotentRequest
	}

	return nil
}

func (m *MockWalletRepository) GetTransactionsHistory(walletID int64) ([]*transaction.Transaction, error) {
	var history []*transaction.Transaction

	t := &transaction.Transaction{ID: 1, FromWalletID: 1, ToWalletID: 2, Amount: 1000,
		CreatedAt: time.Now(), Type: "TRANSFER"}

	history = append(history, t)

	return history, nil
}

func TestWalletService_GetWalletsByProfileID(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	wallets, errWallets := svc.GetWalletsByProfileID(1)

	if errWallets != nil {
		t.Errorf("errWallets: %s", errWallets)
	}

	if len(wallets) == 0 {
		t.Errorf("wallets is empty")
	}
}

func TestWalletService_CreateWallet(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	returnedWalletID, errCreateWallet := svc.CreateWallet(1, "PLN")
	if errCreateWallet != nil {
		t.Errorf("errCreateWallet: %s", errCreateWallet)
	}

	if returnedWalletID != 1 {
		t.Errorf("returnedWalletID is not 1")
	}
}

func TestWalletService_CreateWallet_InvalidCurrency(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	returnedWalletID, errCreateWallet := svc.CreateWallet(1, "AUD")
	if errCreateWallet == nil {
		t.Errorf("Creating wallet went through")
	}

	if returnedWalletID != 0 {
		t.Errorf("returnedWalletID is not 0")
	}

}

func TestWalletService_AddFunds(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	errAddFunds := svc.AddFunds("asdfghjkl123456789", 1, 1, 1000)
	if errAddFunds != nil {
		t.Errorf("errAddFunds is not nil")
	}
}

func TestWalletService_AddFunds_InvalidAmount(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	errAddFunds := svc.AddFunds("asdfghjkl123456789", 1, 1, -1000)
	if errAddFunds == nil {
		t.Errorf("Adding funds went through")
	}
}

func TestWalletService_AddFunds_Idempotency(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	errAddFunds := svc.AddFunds("asdfghjkl1234567890", 1, 1, -1000)
	if errAddFunds != nil {
		t.Errorf("Expected nil error for idempotent request, got: %v", errAddFunds)
	}
}

func TestWalletService_TransferFunds(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	errTransferFunds := svc.TransferFunds("asdfghjkl123456789", 1, 1, 2, 1000)
	if errTransferFunds != nil {
		t.Errorf("errTransferFunds is not nil")
	}
}

func TestWalletService_TransferFunds_InvalidAmount(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	errTransferFunds := svc.TransferFunds("asdfghjkl123456789", 1, 1, 2, -1000)
	if errTransferFunds == nil {
		t.Errorf("Transferring funds went through")
	}
}

func TestWalletService_TransferFunds_SameWallet(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	errTransferFunds := svc.TransferFunds("asdfghjkl123456789", 1, 1, 1, 1000)
	if errTransferFunds == nil {
		t.Errorf("Transferring funds went through")
	}
}

func TestWalletService_TransferFunds_Idempotency(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	errTransferFunds := svc.TransferFunds("asdfghjkl1234567890", 1, 1, 2, 1000)

	if errTransferFunds != nil {
		t.Errorf("Expected nil error for idempotent request, got: %v", errTransferFunds)
	}
}

func TestWalletService_GetTransactionsHistory(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	history, errHistory := svc.GetTransactionsHistory(1, 1)

	if errHistory != nil {
		t.Errorf("errHistory is not nil")
	}

	if len(history) != 1 {
		t.Errorf("history length is not 1")
	}
}

func TestWalletService_GetTransactionsHistory_UnauthorizedAccess(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	svc := NewWalletService(mockRepo)

	history, errHistory := svc.GetTransactionsHistory(1, 10)

	if !errors.Is(errHistory, ErrUnauthorizedAccess) {
		t.Errorf("errHistory is not ErrUnauthorizedAccess")
	}

	if len(history) != 0 {
		t.Errorf("history length is not 0")
	}
}
