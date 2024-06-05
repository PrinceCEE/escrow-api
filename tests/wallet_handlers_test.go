package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/princecee/escrow-api/cmd/app/api/wallets"
	"github.com/princecee/escrow-api/pkg/apis/paystack"
	"github.com/princecee/escrow-api/pkg/json"
	test_utils "github.com/princecee/escrow-api/tests/utils"
	"github.com/princecee/escrow-api/tests/utils/mocks/test_config"
	"github.com/stretchr/testify/suite"
)

type WalletHandlerTestSuite struct {
	suite.Suite
	ts          *test_utils.TestServer
	user        test_utils.TestUser
	accessToken string
}

func (s *WalletHandlerTestSuite) SetupSuite() {
	s.ts = test_utils.NewTestServer()

	user, token := test_utils.SignupPersonalUser(s.ts)
	s.user = user
	s.accessToken = token
}

func (s *WalletHandlerTestSuite) TearDownSuite() {
	s.ts.DropTablesAndTypes()
	s.ts.Config.GetDB().Close()
}

func (s *WalletHandlerTestSuite) get(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", s.accessToken)},
		"Content-Type":  {test_utils.ContentType},
	}

	return req
}

func (s *WalletHandlerTestSuite) post(url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, url, body)
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", s.accessToken)},
		"Content-Type":  {test_utils.ContentType},
	}

	return req
}

func (s *WalletHandlerTestSuite) delete(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", s.accessToken)},
		"Content-Type":  {test_utils.ContentType},
	}

	return req
}

func (s *WalletHandlerTestSuite) TestWalletHandler() {
	url := s.ts.Server.URL + "/api/v1/wallets"
	client := s.ts.Server.Client()

	var bankAccountID, walletID string
	s.Run("manage bank accounts", func() {
		s.Run("add bank account", func() {
			addBankAccountDto := map[string]string{
				"bank_name":      "First Bank",
				"account_name":   "Chimezie Edeh",
				"account_number": "0000000000",
				"bvn":            "00000000000",
			}

			data, _ := json.Marshal(addBankAccountDto)

			req := s.post(url+"/bank-accounts", bytes.NewBuffer(data))
			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[struct {
				BankAccount test_utils.TestBankAccount `json:"bank_account"`
			}])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("account added successfully", respBody.Message)
			s.Equal(addBankAccountDto["account_name"], respBody.Data.BankAccount.AccountName)

			bankAccountID = respBody.Data.BankAccount.ID
			walletID = respBody.Data.BankAccount.WalletID
		})

		s.Run("get bank accounts", func() {
			bankAccounts := []map[string]string{
				{
					"bank_name":      "First Bank",
					"account_name":   "Chimezie Edeh",
					"account_number": "0000000001",
					"bvn":            "00000000001",
				},
				{
					"bank_name":      "First Bank",
					"account_name":   "Chimezie Edeh",
					"account_number": "0000000002",
					"bvn":            "00000000002",
				},
				{
					"bank_name":      "First Bank",
					"account_name":   "Chimezie Edeh",
					"account_number": "0000000003",
					"bvn":            "00000000003",
				},
				{
					"bank_name":      "First Bank",
					"account_name":   "Chimezie Edeh",
					"account_number": "0000000004",
					"bvn":            "00000000004",
				},
				{
					"bank_name":      "First Bank",
					"account_name":   "Chimezie Edeh",
					"account_number": "0000000005",
					"bvn":            "00000000005",
				},
			}

			for _, v := range bankAccounts {
				v := v

				d, _ := json.Marshal(v)
				req := s.post(url+"/bank-accounts", bytes.NewBuffer(d))
				_, err := client.Do(req)
				s.NoError(err)
			}

			req := s.get(fmt.Sprintf("%s/bank-accounts?wallet_id=%s", url, walletID))
			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[struct {
				BankAccounts []test_utils.TestBankAccount `json:"bank_accounts"`
			}])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("bank accounts fetched successfully", respBody.Message)
			s.Len(respBody.Data.BankAccounts, 6)
			s.Equal(6, respBody.Meta.Total)
			s.Equal(1, respBody.Meta.TotalPages)
			s.Equal(1, respBody.Meta.Page)
			s.Equal(20, respBody.Meta.PageSize)
		})

		s.Run("fetch bank accounts with pagination set", func() {
			req := s.get(fmt.Sprintf("%s/bank-accounts?wallet_id=%s&page=1&page_size=3", url, walletID))
			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[struct {
				BankAccounts []test_utils.TestBankAccount `json:"bank_accounts"`
			}])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("bank accounts fetched successfully", respBody.Message)
			s.Len(respBody.Data.BankAccounts, 3)
			s.Equal(6, respBody.Meta.Total)
			s.Equal(2, respBody.Meta.TotalPages)
			s.Equal(1, respBody.Meta.Page)
			s.Equal(3, respBody.Meta.PageSize)
		})

		s.Run("delete bank account", func() {
			req := s.delete(url + fmt.Sprintf("/bank-accounts/%s", bankAccountID))
			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[any])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("bank account deleted successfully", respBody.Message)

			req = s.get(fmt.Sprintf("%s/bank-accounts?wallet_id=%s", url, walletID))
			res, err = client.Do(req)
			s.NoError(err)

			_respBody := new(test_utils.Response[struct {
				BankAccounts []test_utils.TestBankAccount `json:"bank_accounts"`
			}])
			_ = json.ReadJSON(res.Body, _respBody)
			defer res.Body.Close()

			s.Len(_respBody.Data.BankAccounts, 5)
			for _, v := range _respBody.Data.BankAccounts {
				v := v
				s.NotEqual(bankAccountID, v.ID)
			}
		})
	})

	s.Run("manage wallets", func() {
		var ref string
		var walletID string
		fundAmount := 1000000

		s.Run("get wallet", func() {
			req := s.get(url)
			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[struct {
				Wallet test_utils.TestWallet `json:"wallet"`
			}])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("wallet fetched successfully", respBody.Message)
			s.Equal(0, respBody.Data.Wallet.Balance)
			s.Equal(s.user.ID, respBody.Data.Wallet.Identifier)
		})

		s.Run("add funds", func() {
			addFundsDto := map[string]int{
				"amount": fundAmount,
			}

			paystackData := paystack.InitiateTransactionDto{
				Amount: fmt.Sprintf("%d", addFundsDto["amount"]),
				Email:  s.user.Email,
			}
			initiateTransactionResponse := paystack.InitiateTransactionResponse{}
			paystackApi := s.ts.Config.GetAPIs().GetPaystack().(*test_config.TestPaystackAPI)
			paystackApi.On("InitiateTransaction", paystackData).Once().Return(initiateTransactionResponse)

			data, _ := json.Marshal(addFundsDto)
			req := s.post(url+"/add-funds", bytes.NewBuffer(data))

			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[struct {
				WalletHistory test_utils.TestWalletHistory         `json:"wallet_history"`
				PaymentData   paystack.InitiateTransactionResponse `json:"payment_data"`
			}])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("wallet funded successfully", respBody.Message)
			s.Equal(addFundsDto["amount"], respBody.Data.WalletHistory.Amount)
			s.Equal("Pending", respBody.Data.WalletHistory.Status)
			s.Equal(0, respBody.Data.WalletHistory.Wallet.Balance)

			ref = respBody.Data.WalletHistory.ID
			walletID = respBody.Data.WalletHistory.WalletID
		})

		s.Run("get wallet histories", func() {
			req := s.get(url + fmt.Sprintf("/%s/history?status=Pending", walletID))

			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[struct {
				WalletHistories []test_utils.TestWalletHistory `json:"wallet_histories"`
			}])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("wallet histories fetched successfully", respBody.Message)
			s.Len(respBody.Data.WalletHistories, 1)

			walletHistory := respBody.Data.WalletHistories[0]

			s.Equal(ref, walletHistory.ID)
			s.Equal(fundAmount, walletHistory.Amount)
			s.Equal("Pending", walletHistory.Status)
		})

		s.Run("handle webhook", func() {
			webhookDto := new(wallets.WebhookDto[wallets.TransactionData])
			webhookDto.Event = "charge.success"
			webhookDto.Data.Amount = fmt.Sprintf("%d", fundAmount)
			webhookDto.Data.Reference = ref

			data, _ := json.Marshal(webhookDto)
			req := s.post(url+"/paystack-webhook", bytes.NewBuffer(data))

			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[any])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal(http.StatusOK, res.StatusCode)
		})

		s.Run("get wallet histories", func() {
			// fetch pending wallet histories
			req := s.get(url + fmt.Sprintf("/%s/history?status=Pending", walletID))

			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[struct {
				WalletHistories []test_utils.TestWalletHistory `json:"wallet_histories"`
			}])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("wallet histories fetched successfully", respBody.Message)
			s.Len(respBody.Data.WalletHistories, 0)

			// fet successful wallet histories
			req = s.get(url + fmt.Sprintf("/%s/history?status=Successful", walletID))

			res, err = client.Do(req)
			s.NoError(err)

			respBody = new(test_utils.Response[struct {
				WalletHistories []test_utils.TestWalletHistory `json:"wallet_histories"`
			}])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("wallet histories fetched successfully", respBody.Message)
			s.Len(respBody.Data.WalletHistories, 1)

			walletHistory := respBody.Data.WalletHistories[0]

			s.Equal(ref, walletHistory.ID)
			s.Equal(fundAmount, walletHistory.Amount)
			s.Equal("Successful", walletHistory.Status)
			s.Equal(fundAmount, walletHistory.Wallet.Balance)
		})

		s.Run("withdraw funds", func() {
			req := s.get(fmt.Sprintf("%s/bank-accounts?wallet_id=%s&page=1&page_size=1", url, walletID))
			res, err := client.Do(req)
			s.NoError(err)

			respBody := new(test_utils.Response[struct {
				BankAccounts []test_utils.TestBankAccount `json:"bank_accounts"`
			}])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			bankAccount := respBody.Data.BankAccounts[0]

			withdrawFundDto := map[string]any{
				"amount":          2000000,
				"bank_account_id": bankAccount.ID,
			}

			s.Run("withddraw more than available", func() {
				data, _ := json.Marshal(withdrawFundDto)

				req := s.post(url+"/withdraw-funds", bytes.NewBuffer(data))
				res, err := client.Do(req)
				s.NoError(err)

				respBody := new(test_utils.Response[any])
				_ = json.ReadJSON(res.Body, respBody)
				defer res.Body.Close()

				s.Equal(false, respBody.Success)
				s.Equal("insufficient balance", respBody.Message)
			})

			s.Run("withdraw available", func() {
				withdrawFundDto["amount"] = 500000
				data, _ := json.Marshal(withdrawFundDto)

				req := s.post(url+"/withdraw-funds", bytes.NewBuffer(data))
				res, err := client.Do(req)
				s.NoError(err)

				respBody := new(test_utils.Response[any])
				_ = json.ReadJSON(res.Body, respBody)
				defer res.Body.Close()

				s.Equal(true, respBody.Success)
				s.Equal("funds successfully withdrawn", respBody.Message)
			})
		})
	})
}

func TestWalletHandlerSuite(t *testing.T) {
	suite.Run(t, &WalletHandlerTestSuite{})
}
