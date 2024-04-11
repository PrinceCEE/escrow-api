package models

type Wallet struct {
	Balance     int       `json:"balance" db:"balance"`
	Receivable  int       `json:"receivable_balance" db:"receivable_balance"`
	Payable     int       `json:"payable_balance" db:"payable_balance"`
	AccountType string    `json:"account_type" db:"account_type"`
	Identifier  string    `json:"identifier" db:"identifier"`
	User        *User     `json:"user,omitempty" db:"-"`
	Business    *Business `json:"business,omitempty" db:"-"`
	ModelMixin
}
