package wallet

type Wallet struct {
	ID        int64  `json:"id"`
	ProfileID int64  `json:"profile_id"`
	Balance   int64  `json:"balance"`
	Currency  string `json:"currency"`
}
