package wallets

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/json"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/Bupher-Co/bupher-api/pkg/validator"
	"github.com/jackc/pgx/v5"
)

type walletHandler struct {
	c config.IConfig
}

func (h *walletHandler) addBankAccount(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(addNewAccountDto)

	tx, _ := h.c.GetDB().Begin(context.Background())
	defer tx.Rollback(context.Background())

	walletRepo := h.c.GetWalletRepository()
	bankAccountRepo := h.c.GetBankAccountRepository()
	businessRepo := h.c.GetBusinessRepository()

	err := json.ReadJSON(r.Body, body)
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

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, err = walletRepo.GetByIdentifier(user.ID.String(), tx)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				resp.Message = err.Error()
				response.SendErrorResponse(w, resp, http.StatusUnauthorized)
				return
			}

			wallet = &models.Wallet{AccountType: user.AccountType, Identifier: user.ID}
			err = walletRepo.Create(wallet, tx)
			if err != nil {
				resp.Message = err.Error()
				response.SendErrorResponse(w, resp, http.StatusInternalServerError)
				return
			}
		}
	} else {
		business, err := businessRepo.GetByUserID(user.ID.String(), tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		wallet, err = walletRepo.GetByIdentifier(business.ID.String(), tx)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				resp.Message = err.Error()
				response.SendErrorResponse(w, resp, http.StatusUnauthorized)
				return
			}

			wallet = &models.Wallet{AccountType: user.AccountType, Identifier: business.ID}
			err = walletRepo.Create(wallet, tx)
			if err != nil {
				resp.Message = err.Error()
				response.SendErrorResponse(w, resp, http.StatusInternalServerError)
				return
			}
		}
	}

	bankAccount := &models.BankAccount{
		BankName:      body.BankName,
		AccountName:   body.AccountName,
		AccountNumber: body.AccountNumber,
		BVN:           body.BVN,
		WalletID:      wallet.ID,
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
	businessRepo := h.c.GetBusinessRepository()

	user := r.Context().Value(utils.ContextKey{}).(*models.User)
	business, _ := businessRepo.GetByUserID(user.ID.String(), nil)
	bankAccountId := r.URL.Query().Get("bank_account_id")

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
		wallet, _ = walletRepo.GetByIdentifier(user.ID.String(), nil)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(business.ID.String(), nil)
	}

	if bankAccount.WalletID.String() != wallet.ID.String() {
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
	businessRepo := h.c.GetBusinessRepository()

	user := r.Context().Value(utils.ContextKey{}).(*models.User)
	business, _ := businessRepo.GetByUserID(user.ID.String(), nil)
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
		wallet, _ = walletRepo.GetByIdentifier(user.ID.String(), nil)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(business.ID.String(), nil)
	}

	if wallet.ID.String() != body.WalletID {
		resp.Message = "forbidden"
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	offset := (body.Page - 1) * body.PageSize
	args := []any{body.WalletID, offset, body.PageSize}
	rows, err := h.c.GetDB().Query(context.Background(), `
		SELECT * FROM bank_accounts
		WHERE wallet_id = $1
		ORDER BY created_at DESC
		OFFSET $2
		LIMIT $3`,
		args...,
	)
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

	bankAccounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.BankAccount])
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Message = "bank accounts fetched successfully"
	resp.Data = map[string][]models.BankAccount{
		"bank_accounts": bankAccounts,
	}
	response.SendResponse(w, resp)
}

func (h *walletHandler) getWallet(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	walletRepo := h.c.GetWalletRepository()
	businessRepo := h.c.GetBusinessRepository()

	user := r.Context().Value(utils.ContextKey{}).(*models.User)
	business, _ := businessRepo.GetByUserID(user.ID.String(), nil)

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, _ = walletRepo.GetByIdentifier(user.ID.String(), nil)
	} else {
		wallet, _ = walletRepo.GetByIdentifier(business.ID.String(), nil)
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
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	tx, _ := h.c.GetDB().Begin(context.Background())

	walletRepo := h.c.GetWalletRepository()
	walletHistoryRepo := h.c.GetWalletHistoryRepository()
	businessRepo := h.c.GetBusinessRepository()

	user := r.Context().Value(utils.ContextKey{}).(*models.User)

	wallet := new(models.Wallet)
	if user.AccountType == models.PersonalAccountType {
		wallet, _ = walletRepo.GetByIdentifier(user.ID.String(), tx)
	} else {
		business, _ := businessRepo.GetByUserID(user.ID.String(), tx)
		wallet, _ = walletRepo.GetByIdentifier(business.ID.String(), tx)
	}

	wallet.Balance += body.Amount
	wallet.Receivable += body.Amount
	err = walletRepo.Update(wallet, tx)
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
	}
	err = walletHistoryRepo.Create(walletHistory, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	// read into how to populate the struct in a model
	// change the dependency to [user depends on business]
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) withrawFunds(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) getWalletHistories(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}
