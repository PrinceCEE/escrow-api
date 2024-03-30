package wallets

type addNewAccountDto struct {
	BankName      string `json:"bank_name" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required,len=10"`
	BVN           string `json:"bvn" validate:"required,len=11"`
}

type getBankAccountsQueryDto struct {
	WalletID string `json:"wallet_id" validate:"uuid"`
	Page     int    `json:"page" validate:"number,min=1"`
	PageSize int    `json:"page_size" validate:"number,min=1,max=100"`
}

type addFundsDto struct {
	Amount int `json:"amount" validate:"required,min=5000"`
}

type withrawFundsDto struct {
	Amount        int    `json:"amount" validate:"required,min=5000"`
	BankAccountId string `json:"bank_account_id" validate:"required,uuid"`
}

type getWalletHistoriesQueryDto getBankAccountsQueryDto
