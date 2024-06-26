package wallets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/princecee/escrow-api/cmd/app/pkg/response"
	"github.com/princecee/escrow-api/config"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/pkg/apis/paystack"
	"github.com/princecee/escrow-api/pkg/json"
	"github.com/princecee/escrow-api/pkg/utils"
	"github.com/princecee/escrow-api/pkg/validator"
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
		wallet, err = walletRepo.GetByIdentifier(*user.BusinessID, tx)
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
		Wallet:        wallet,
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
		wallet, _ = walletRepo.GetByIdentifier(*user.BusinessID, nil)
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
		wallet, _ = walletRepo.GetByIdentifier(*user.BusinessID, nil)
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
	_ = h.c.GetDB().
		QueryRow(
			context.Background(),
			`SELECT COUNT(*) AS total FROM bank_accounts b
			INNER JOIN wallets w ON w.id = b.wallet_id
			WHERE b.wallet_id = $1`,
			body.WalletID,
		).
		Scan(&total)

	resp.Meta.Page = body.Page
	resp.Meta.PageSize = body.PageSize
	resp.Meta.Total = total
	resp.Meta.TotalPages = utils.GetTotalPages(total, body.PageSize)

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
		wallet, _ = walletRepo.GetByIdentifier(*user.BusinessID, nil)
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
	paystackAPI := h.c.GetAPIs().GetPaystack()

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
		wallet, _ = walletRepo.GetByIdentifier(*user.BusinessID, tx)
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

	paystackResponse, err := paystackAPI.InitiateTransaction(paystack.InitiateTransactionDto{
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
	response.SendResponse(w, resp)
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
	defer tx.Rollback(context.Background())

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, _ = walletRepo.GetByIdentifier(user.ID, tx)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(*user.BusinessID, tx)
	}

	if wallet.Receivable-body.Amount < 0 {
		resp.Message = "insufficient balance"
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
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

	status := query.Get("status")
	wType := query.Get("type")
	walletId := chi.URLParam(r, "wallet_id")

	body := &getWalletHistoriesQueryDto{
		WalletID: walletId,
		Page:     utils.GetPage(int(page)),
		PageSize: utils.GetPageSize(int(pageSize)),
	}

	if status != "" {
		body.Status = status
	}
	if wType != "" {
		body.Type = wType
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
		wallet, _ = walletRepo.GetByIdentifier(*user.BusinessID, nil)
	}

	if body.WalletID != wallet.ID {
		resp.Message = "forbidden"
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	pagination := utils.GetPagination(body.Page, body.PageSize)
	where, args := utils.GenerateANDWhereFromArgs([]utils.WhereArgs{
		{
			Name:  "h.wallet_id",
			Value: wallet.ID,
		},
		{
			Name:  "h.status",
			Value: body.Status,
		},
		{
			Name:  "h.type",
			Value: body.Type,
		},
	})

	var total int
	_ = h.c.GetDB().
		QueryRow(
			context.Background(),
			fmt.Sprintf(
				`SELECT COUNT(*) FROM wallet_histories h
				INNER JOIN wallets w ON w.id = h.wallet_id
				%s
				`,
				where,
			),
			args...,
		).
		Scan(&total)

	args = append(args, pagination.Offset, pagination.Limit)
	walletHistories, err := walletHistoryRepo.GetMany(args, where, nil)

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

	resp.Message = "wallet histories fetched successfully"
	resp.Data = map[string]any{
		"wallet_histories": walletHistories,
	}
	resp.Meta.Page = body.Page
	resp.Meta.PageSize = body.PageSize
	resp.Meta.Total = total
	resp.Meta.TotalPages = utils.GetTotalPages(total, body.PageSize)

	response.SendResponse(w, resp)
}

func (h *walletHandler) handlePaystackWebhook(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}

	env := h.c.Getenv("ENVIRONMENT")

	if env != "test" {
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
	}

	tx, _ := h.c.GetDB().Begin(context.Background())
	defer tx.Rollback(context.Background())

	// read the untyped response body so as to know the event
	tmp := make(map[string]any)
	err := json.ReadJSON(r.Body, &tmp)
	defer r.Body.Close()

	if err != nil {
		h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	switch tmp["event"] {
	case "charge.success":
		body := new(WebhookDto[TransactionData])

		tmpJsonByte, _ := json.Marshal(tmp)
		err = json.Unmarshal(tmpJsonByte, body)

		if err != nil {
			h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		walletHistoryRepo := h.c.GetWalletHistoryRepository()
		walletRepo := h.c.GetWalletRepository()
		transactionRepo := h.c.GetTransactionRepository()
		transactionTimelineRepo := h.c.GetTransactionTimelineRepository()

		isForTransaction, ok := body.Data.Metadata.(map[string]any)["is_for_transaction"]
		if !ok {
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
		} else {
			if !isForTransaction.(bool) {
				if err != nil {
					h.c.GetLogger().Log(zerolog.InfoLevel, "is_for_transaction wrongly placed", nil, err)
					response.SendErrorResponse(w, resp, http.StatusBadRequest)
					return
				}
			}

			transaction, err := transactionRepo.GetById(body.Data.Reference, tx)
			if err != nil {
				h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
				response.SendErrorResponse(w, resp, http.StatusBadRequest)
				return
			}

			transaction.Status = models.TransactionStatusPendingDelivery
			err = transactionRepo.Update(transaction, tx)
			if err != nil {
				h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
				response.SendErrorResponse(w, resp, http.StatusBadRequest)
				return
			}

			timeline := &models.TransactionTimeline{
				TransactionID: transaction.ID,
				Name:          models.TimelinePaymentSubmitted,
			}
			err = transactionTimelineRepo.Create(timeline, tx)
			if err != nil {
				h.c.GetLogger().Log(zerolog.InfoLevel, err.Error(), nil, err)
				response.SendErrorResponse(w, resp, http.StatusBadRequest)
				return
			}
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
