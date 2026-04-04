package transaction

import "time"

type Transaction struct {
	ID           int64
	FromWalletID int64
	ToWalletID   int64
	Amount       int64
	CreatedAt    time.Time
}
