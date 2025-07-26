package accounts

import "time"

type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"` // in cents
	CreatedAt time.Time `json:"created_at"`
}
