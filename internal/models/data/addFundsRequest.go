package data

type AddFundsRequest struct {
	WalletID int64 `json:"wallet_id"`
	Amount   int64 `json:"amount"`
}
