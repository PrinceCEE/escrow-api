package wallets

import "time"

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

type tranactionData struct {
	ID              int       `json:"id"`
	Domain          string    `json:"domain"`
	Status          string    `json:"status"`
	Reference       string    `json:"reference"`
	Amount          string    `json:"amount"`
	Message         string    `json:"message"`
	GatewayResponse string    `json:"gateway_response"`
	PaidAt          time.Time `json:"paid_at"`
	CreatedAt       time.Time `json:"created_at"`
	Channel         string    `json:"channel"`
	Currency        string    `json:"currency"`
	IpAddress       string    `json:"ip_address"`
	Metadata        any       `json:"metadata"`
	Log             struct {
		TimeSpent      int    `json:"time_spent"`
		Attempts       int    `json:"attempts"`
		Authentication string `json:"authentication"`
		Errors         string `json:"errors"`
		Success        bool   `json:"success"`
		Mobile         string `json:"mobile"`
		Input          any    `json:"input"`
		Channel        any    `json:"channel"`
		History        []struct {
			Input   string `json:"type"`
			Message string `json:"message"`
			Time    int    `json:"time"`
		}
	}
	Fees     any `json:"fees"`
	Customer struct {
		ID           int    `json:"id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		CustomerCode string `json:"customer_code"`
		Phone        string `json:"phone"`
		MetaData     any    `json:"metadata"`
		RiskAction   string `json:"risk_action"`
	}
	Authorization struct {
		AuthorizationCode string `json:"authorization_code"`
		Bin               string `json:"bin"`
		Last4             string `json:"last4"`
		ExpMonth          string `json:"exp_month"`
		ExpYear           string `json:"exp_year"`
		CardType          string `json:"card_type"`
		Bank              string `json:"bank"`
		CountryCode       string `json:"country_code"`
		Brand             string `json:"brand"`
		AccountName       string `json:"account_name"`
	}
	Plan any `json:"plan"`
}

type webhookDto[T any] struct {
	Event string `json:"event"`
	Data  T      `json:"data"`
}
