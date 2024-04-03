package models

type BankAccount struct {
	BankName      string `json:"bank_name" db:"bank_name"`
	AccountName   string `json:"account_name" db:"account_name"`
	AccountNumber string `json:"account_number" db:"account_number"`
	BVN           string `json:"bvn" db:"bvn"`
	WalletID      string `json:"wallet_id" db:"wallet_id"`
	Wallet        Wallet `json:"wallet,omitempty" db:"-"`
	ModelMixin
}
