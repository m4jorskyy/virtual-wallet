package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"virtual-wallet/internal/models/data"
	"virtual-wallet/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (s *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	var request data.RegisterRequest
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&request)

	if errDecode != nil {
		http.Error(w, fmt.Sprintf("Error decoding: %s", errDecode), http.StatusBadRequest)
		return
	}

	returnedID, errRegister := s.service.RegisterNewUser(request.FirstName, request.LastName, request.Email,
		request.Login, request.Password)

	if errRegister != nil {
		http.Error(w, fmt.Sprintf("Error registering: %s", errRegister), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	errEncoding := encoder.Encode(map[string]int64{"id": returnedID})

	if errEncoding != nil {
		http.Error(w, fmt.Sprintf("Error encoding: %s", errEncoding), http.StatusInternalServerError)
		return
	}
}
