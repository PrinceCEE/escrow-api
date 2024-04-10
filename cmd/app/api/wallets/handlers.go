package wallets

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/apis/paystack"
	"github.com/Bupher-Co/bupher-api/pkg/json"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/Bupher-Co/bupher-api/pkg/validator"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type walletHandler struct {
	c config.IConfig
}

func (h *walletHandler) addBankAccount(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(addNewAccountDto)

	walletRepo := h.c.GetWalletRepository()
	bankAccountRepo := h.c.GetBankAccountRepository()

	err := json.ReadJSON(r.Body, body)
	defer r.Body.Close()

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		resp.Message = response.ErrBadRequest.Error()
		resp.Data = validationErrors
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	tx, _ := h.c.GetDB().Begin(context.Background())
	defer tx.Rollback(context.Background())

	user := r.Context().Value(utils.ContextKey{}).(*models.User)

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, err = walletRepo.GetByIdentifier(user.ID, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}
	} else {
		wallet, err = walletRepo.GetByIdentifier(user.BusinessID, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}
	}

	bankAccount := &models.BankAccount{
		BankName:      body.BankName,
		AccountName:   body.AccountName,
		AccountNumber: body.AccountNumber,
		BVN:           body.BVN,
		WalletID:      wallet.ID,
		Wallet:        *wallet,
	}
	err = bankAccountRepo.Create(bankAccount, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Message = "account added successfully"
	resp.Data = map[string]models.BankAccount{
		"bank_account": *bankAccount,
	}
	response.SendResponse(w, resp)
}

func (h *walletHandler) deleteBankAccount(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	bankAccountRepo := h.c.GetBankAccountRepository()
	walletRepo := h.c.GetWalletRepository()

	user := r.Context().Value(utils.ContextKey{}).(*models.User)
	bankAccountId := chi.URLParam(r, "bank_account_id")

	bankAccount, err := bankAccountRepo.GetById(bankAccountId, nil)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = "bank account not found"
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, _ = walletRepo.GetByIdentifier(user.ID, nil)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(user.BusinessID, nil)
	}

	if bankAccount.WalletID != wallet.ID {
		resp.Message = "forbidden"
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	err = bankAccountRepo.Delete(bankAccountId, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	resp.Message = "bank account deleted successfully"
	response.SendResponse(w, resp)
}

func (h *walletHandler) getBankAccounts(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	walletRepo := h.c.GetWalletRepository()
	bankAccountRepo := h.c.GetBankAccountRepository()

	user := r.Context().Value(utils.ContextKey{}).(*models.User)
	query := r.URL.Query()

	var page, pageSize int64
	if query.Get("page") != "" {
		page, _ = strconv.ParseInt(query.Get("page"), 10, 64)
	}
	if query.Get("page_size") != "" {
		pageSize, _ = strconv.ParseInt(query.Get("page_size"), 10, 64)
	}

	body := &getBankAccountsQueryDto{
		WalletID: query.Get("wallet_id"),
		Page:     utils.GetPage(int(page)),
		PageSize: utils.GetPageSize(int(pageSize)),
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		resp.Message = response.ErrBadRequest.Error()
		resp.Data = validationErrors
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, _ = walletRepo.GetByIdentifier(user.ID, nil)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(user.BusinessID, nil)
	}

	if wallet.ID != body.WalletID {
		resp.Message = "forbidden"
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	pagination := utils.GetPagination(body.Page, body.PageSize)

	bankAccounts, err := bankAccountRepo.GetByWalletId(body.WalletID, pagination, nil)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = "bank accounts fetched successfully"
			resp.Data = map[string][]models.BankAccount{
				"bank_accounts": {},
			}
			response.SendResponse(w, resp)
			return
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	var total int
	_ = h.c.GetDB().QueryRow(context.Background(), `
		SELECT COUNT(*) AS total FROM bank_accounts b
		INNER JOIN wallets w ON w.id = b.wallet_id
		WHERE b.wallet_id = $1
	`, body.WalletID).Scan(&total)

	resp.Meta.Page = body.Page
	resp.Meta.PageSize = body.PageSize
	resp.Meta.Total = total
	resp.Meta.TotalPages = int(math.Ceil((float64(total) / float64(body.PageSize))))

	resp.Message = "bank accounts fetched successfully"
	resp.Data = map[string][]*models.BankAccount{
		"bank_accounts": bankAccounts,
	}

	response.SendResponse(w, resp)
}

func (h *walletHandler) getWallet(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	walletRepo := h.c.GetWalletRepository()

	user := r.Context().Value(utils.ContextKey{}).(*models.User)

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, _ = walletRepo.GetByIdentifier(user.ID, nil)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(user.BusinessID, nil)
	}

	resp.Message = "wallet fetched successfully"
	resp.Data = map[string]models.Wallet{
		"wallet": *wallet,
	}
	response.SendResponse(w, resp)
}

func (h *walletHandler) addFunds(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(addFundsDto)

	err := json.ReadJSON(r.Body, body)
	defer r.Body.Close()

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		resp.Message = response.ErrBadRequest.Error()
		resp.Data = validationErrors
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	tx, _ := h.c.GetDB().Begin(context.Background())

	walletRepo := h.c.GetWalletRepository()
	walletHistoryRepo := h.c.GetWalletHistoryRepository()

	user := r.Context().Value(utils.ContextKey{}).(*models.User)

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, _ = walletRepo.GetByIdentifier(user.ID, tx)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(user.BusinessID, tx)
	}

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	walletHistory := &models.WalletHistory{
		WalletID: wallet.ID,
		Type:     models.WalletHistoryDepositType,
		Amount:   body.Amount,
		Status:   models.WalletHistoryPending,
		Wallet:   *wallet,
	}
	err = walletHistoryRepo.Create(walletHistory, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	paystackResponse, err := paystack.InitiateTransaction(paystack.InitiateTransactionDto{
		Email:     user.Email,
		Amount:    strconv.FormatInt(int64(body.Amount), 10),
		Reference: walletHistory.ID,
	})
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Message = "wallet funded successfully"
	resp.Data = map[string]any{
		"wallet_history": walletHistory,
		"payment_data":   paystackResponse,
	}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) withrawFunds(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(withrawFundsDto)

	err := json.ReadJSON(r.Body, body)
	defer r.Body.Close()

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		resp.Message = response.ErrBadRequest.Error()
		resp.Data = validationErrors
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	user := r.Context().Value(utils.ContextKey{}).(*models.User)

	walletRepo := h.c.GetWalletRepository()
	walletHistoryRepo := h.c.GetWalletHistoryRepository()

	tx, _ := h.c.GetDB().Begin(context.Background())

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, _ = walletRepo.GetByIdentifier(user.ID, tx)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(user.BusinessID, tx)
	}

	wallet.Balance -= body.Amount
	wallet.Receivable -= body.Amount
	err = walletRepo.Update(wallet, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	walletHistory := &models.WalletHistory{
		WalletID: wallet.ID,
		Type:     models.WalletHistoryWithdrawalType,
		Amount:   body.Amount,
		Status:   models.WalletHistoryPending,
		Wallet:   *wallet,
	}
	err = walletHistoryRepo.Create(walletHistory, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	// TODO - make transfer through paystack

	resp.Message = "funds successfully withdrawn"
	response.SendResponse(w, resp)
}

func (h *walletHandler) getWalletHistories(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}

	walletRepo := h.c.GetWalletRepository()
	walletHistoryRepo := h.c.GetWalletHistoryRepository()

	user := r.Context().Value(utils.ContextKey{}).(*models.User)
	query := r.URL.Query()

	var page, pageSize int64
	if query.Get("page") != "" {
		page, _ = strconv.ParseInt(query.Get("page"), 10, 64)
	}
	if query.Get("page_size") != "" {
		pageSize, _ = strconv.ParseInt(query.Get("page_size"), 10, 64)
	}

	body := &getWalletHistoriesQueryDto{
		WalletID: query.Get("wallet_id"),
		Page:     utils.GetPage(int(page)),
		PageSize: utils.GetPageSize(int(pageSize)),
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		resp.Message = response.ErrBadRequest.Error()
		resp.Data = validationErrors
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, _ = walletRepo.GetByIdentifier(user.ID, nil)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(user.BusinessID, nil)
	}

	if body.WalletID != wallet.ID {
		resp.Message = "forbidden"
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	pagination := utils.GetPagination(body.Page, body.PageSize)
	walletHistories, err := walletHistoryRepo.GetByWalletId(body.WalletID, pagination, nil)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = "wallet histories fetched successfully"
			resp.Data = map[string]any{
				"wallet_histories": []any{},
			}
			response.SendResponse(w, resp)
			return
		default:
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}
	}

	var total int
	_ = h.c.GetDB().QueryRow(context.Background(), `
		SELECT COUNT(*) AS total FROM wallet_histories h
		INNER JOIN wallets w ON w.id = h.wallet_id
		WHERE h.wallet_id = $1
	`, body.WalletID).Scan(&total)

	resp.Message = "wallet histories fetched successfully"
	resp.Data = map[string]any{
		"wallet_histories": walletHistories,
	}
	resp.Meta.Page = body.Page
	resp.Meta.PageSize = body.PageSize
	resp.Meta.Total = total
	resp.Meta.TotalPages = int(total / body.PageSize)

	response.SendResponse(w, resp)
}

func (h *walletHandler) handlePaystackWebhook(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}

	paystackSig := r.Header.Get("x-paystack-signature")
	if paystackSig == "" {
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	hash, err := utils.ComputeHMAC(r.Body)
	if err != nil {
		h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	if hash != paystackSig {
		h.c.GetLogger().Log(zerolog.InfoLevel, "forbidden", nil, err)
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	tx, _ := h.c.GetDB().Begin(context.Background())

	// read the untyped response body so as to know the event
	var tmp map[string]any
	err = json.ReadJSON(r.Body, tmp)
	if err != nil {
		h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	switch tmp["event"] {
	case "charge.success":
		body, err := json.ReadTypedJSON[webhookDto[tranactionData]](r.Body)
		defer r.Body.Close()

		if err != nil {
			h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		validationErrors := validator.ValidateData(body)
		if validationErrors != nil {
			h.c.GetLogger().Log(zerolog.InfoLevel, "validation error", validationErrors, err)
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		walletHistoryRepo := h.c.GetWalletHistoryRepository()
		walletRepo := h.c.GetWalletRepository()

		walletHistory, err := walletHistoryRepo.GetById(body.Data.Reference, tx)
		if err != nil {
			h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		wallet, err := walletRepo.GetById(walletHistory.WalletID, tx)
		if err != nil {
			h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		walletHistory.Status = models.WalletHistorySuccessful
		wallet.Balance += walletHistory.Amount
		wallet.Receivable += walletHistory.Amount

		err = walletHistoryRepo.Update(walletHistory, tx)
		if err != nil {
			h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		err = walletRepo.Update(wallet, tx)
		if err != nil {
			h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}
	response.SendResponse(w, resp)
}
