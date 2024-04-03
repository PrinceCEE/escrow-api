package models

const (
	WalletHistoryWithdrawalType = "Withdrawal"
	WalletHistoryDepositType    = "Deposit"
)

const (
	WalletHistorySuccessful = "Successful"
	WalletHistoryCanceled   = "Canceled"
	WalletHistoryPending    = "Pending"
)

type WalletHistory struct {
	WalletID string `json:"wallet_id" db:"wallet_id"`
	Type     string `json:"type" db:"type"`
	Amount   int    `json:"amount" db:"amount"`
	Status   string `json:"status" db:"status"`
	Wallet   Wallet `json:"wallet,omitempty" db:"-"`
	ModelMixin
}
