package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"virtual-wallet/internal/models/data"
	"virtual-wallet/internal/service"
)

type WalletHandler struct {
	service *service.WalletService
}

func NewWalletHandler(service *service.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (s *WalletHandler) GetWalletsByProfileID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	profileID, ok := r.Context().Value(userContextKey).(int64)

	if !ok {
		http.Error(w, "Invalid userID", http.StatusInternalServerError)
		return
	}

	wallets, errWallets := s.service.GetWalletsByProfileID(profileID)

	if errWallets != nil {
		http.Error(w, "Error getting wallets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	errEncoding := encoder.Encode(&wallets)

	if errEncoding != nil {
		http.Error(w, fmt.Sprintf("Error encoding: %s", errEncoding), http.StatusInternalServerError)
		return
	}
}

func (s *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	profileID, ok := r.Context().Value(userContextKey).(int64)

	if !ok {
		http.Error(w, "Invalid userID", http.StatusInternalServerError)
		return
	}

	var request data.CreateWalletRequest
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&request)

	if errDecode != nil {
		http.Error(w, fmt.Sprintf("Error decoding: %s", errDecode), http.StatusBadRequest)
		return
	}

	returnedWalletID, errReturnedWalletID := s.service.CreateWallet(profileID, request.Currency)

	if errReturnedWalletID != nil {
		http.Error(w, fmt.Sprintf("Error getting walletID: %s", errReturnedWalletID), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	errEncode := encoder.Encode(map[string]int64{"id": returnedWalletID})

	if errEncode != nil {
		http.Error(w, fmt.Sprintf("Error encoding: %s", errEncode), http.StatusInternalServerError)
	}

}

func (s *WalletHandler) AddFunds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	profileID, ok := r.Context().Value(userContextKey).(int64)

	if !ok {
		http.Error(w, "Invalid userID", http.StatusInternalServerError)
		return
	}

	var request data.AddFundsRequest
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&request)

	if errDecode != nil {
		http.Error(w, fmt.Sprintf("Error decoding: %s", errDecode), http.StatusBadRequest)
		return
	}

	errAddFunds := s.service.AddFunds(request.WalletID, profileID, request.Amount)

	if errAddFunds != nil {
		if errors.Is(errAddFunds, service.ErrInvalidAmount) {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		http.Error(w, fmt.Sprintf("Error adding funds: %s", errAddFunds), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
