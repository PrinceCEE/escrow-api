package wallets

import (
	"context"
	"errors"
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/json"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/jackc/pgx/v5"
)

type walletHandler struct {
	c config.IConfig
}

func (h *walletHandler) addBankAccount(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(addNewAccountDto)

	tx, err := h.c.GetDB().Begin(context.Background())
	defer tx.Rollback(context.Background())

	walletRepo := h.c.GetWalletRepository()
	bankAccountRepo := h.c.GetBankAccountRepository()
	businessRepo := h.c.GetBusinessRepository()

	err = json.ReadJSON(r.Body, body)
	if err != nil {
		resp.Message = err.Error()
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

		wallet, err := walletRepo.GetByIdentifier(business.ID.String(), tx)
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
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) getBankAccounts(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) addFunds(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) withrawFunds(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) getWallet(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) getWalletHistories(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}
