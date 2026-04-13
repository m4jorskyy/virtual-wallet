package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"virtual-wallet/internal/models/wallet"
	"virtual-wallet/internal/service"
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

func (m *MockWalletRepository) AddFunds(walletID int64, profileID int64, amount int64) error {
	return nil
}

func TestWalletHandler_GetWalletsByProfileID(t *testing.T) {
	request := httptest.NewRequest("GET", "/api/wallets/", strings.NewReader(""))
	ctx := context.WithValue(request.Context(), userContextKey, int64(1))
	request = request.WithContext(ctx)

	recorder := httptest.NewRecorder()
	mockRepo := &MockWalletRepository{}
	svc := service.NewWalletService(mockRepo)
	handler := NewWalletHandler(svc)
	handler.GetWalletsByProfileID(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetWalletsByProfileID code is not 200")
	}
}

func TestWalletHandler_CreateWallet(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/wallet/create", strings.NewReader("{\n	\"currency\": \"PLN\"\n}"))
	ctx := context.WithValue(request.Context(), userContextKey, int64(1))
	request = request.WithContext(ctx)

	recorder := httptest.NewRecorder()
	mockRepo := &MockWalletRepository{}
	svc := service.NewWalletService(mockRepo)
	handler := NewWalletHandler(svc)
	handler.CreateWallet(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("CreateWallet code is not 201")
	}
}

func TestNewWalletHandler_AddFunds(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/wallet/addFunds", strings.NewReader("{\n	\"wallet_id\": 1,	\n	\"amount\": 1000	\n}"))
	ctx := context.WithValue(request.Context(), userContextKey, int64(1))
	request = request.WithContext(ctx)

	recorder := httptest.NewRecorder()
	mockRepo := &MockWalletRepository{}
	svc := service.NewWalletService(mockRepo)
	handler := NewWalletHandler(svc)
	handler.AddFunds(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("AddFunds code is not 200")
	}
}
