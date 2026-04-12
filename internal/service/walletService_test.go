package service

import (
	"testing"
	"virtual-wallet/internal/models/wallet"
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
