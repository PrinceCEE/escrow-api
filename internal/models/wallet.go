package models

import uuid "github.com/jackc/pgx-gofrs-uuid"

type Wallet struct {
	Balance     int       `json:"balance" db:"balance"`
	Receivable  int       `json:"receivable_balance" db:"receivable_balance"`
	Payable     int       `json:"payable_balance" db:"payable_balance"`
	AccountType string    `json:"account_type" db:"account_type"`
	Identifier  uuid.UUID `json:"identifier" db:"identifier"`
	ModelMixin
}