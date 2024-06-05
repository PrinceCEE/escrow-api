package users

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/princecee/escrow-api/cmd/app/pkg/response"
	"github.com/princecee/escrow-api/config"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/pkg/json"
	"github.com/princecee/escrow-api/pkg/utils"
	"github.com/princecee/escrow-api/pkg/validator"
)

type userHandler struct {
	c config.IConfig
}

func (h *userHandler) getMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(utils.ContextKey{}).(*models.User)

	resp := response.ApiResponse{Message: "user fetched successfully", Data: map[string]any{
		"user": user,
	}}

	response.SendResponse(w, resp)
}

func (h *userHandler) getUser(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}

	userId := chi.URLParam(r, "user_id")
	userRepo := h.c.GetUserRepository()

	user, err := userRepo.GetById(userId, nil)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = "user not found"
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusNotFound)
		return
	}

	resp.Message = "user fetched successfully"
	resp.Data = map[string]any{
		"user": user,
	}
	response.SendResponse(w, resp)
}

func (h *userHandler) updateAccount(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}

	body := new(updateAccountDto)
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

	userRepo := h.c.GetUserRepository()
	businessRepo := h.c.GetBusinessRepository()

	if body.FirstName != nil {
		user.FirstName = models.NullString{NullString: sql.NullString{String: *body.FirstName, Valid: true}}
	}
	if body.LastName != nil {
		user.LastName = models.NullString{NullString: sql.NullString{String: *body.LastName, Valid: true}}
	}
	if body.PhoneNumber != nil {
		user.PhoneNumber = models.NullString{NullString: sql.NullString{String: *body.PhoneNumber, Valid: true}}
	}
	if body.Email != nil {
		user.Email = *body.Email
	}
	if body.ImageUrl != nil {
		user.ImageUrl = *body.ImageUrl
	}

	err = userRepo.Update(user, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	if user.AccountType == models.BusinessAccountType {
		business, err := businessRepo.GetById(*user.BusinessID, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		if body.Email != nil {
			business.Email = *body.Email
		}
		if body.BusinessName != nil {
			business.Name = *body.BusinessName
		}
		if body.ImageUrl != nil {
			business.ImageUrl = *body.ImageUrl
		}

		err = businessRepo.Update(business, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		user.Business = business
	}

	err = tx.Commit(context.Background())
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Data = map[string]any{"user": user}
	resp.Message = "account updated successfully"
	response.SendResponse(w, resp)
}

func (h *userHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}

	body := new(changePasswordDto)
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
	authRepo := h.c.GetAuthRepository()

	password, _ := utils.GeneratePasswordHash(body.Password)
	auth, err := authRepo.GetByUserId(user.ID, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	oldPassword := auth.Password
	auth.PasswordHistory = append(auth.PasswordHistory, models.PasswordHistory{
		Password:  oldPassword,
		Timestamp: time.Now(),
	})
	auth.Password = string(password)

	err = authRepo.Update(auth, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Message = "password changed successfully"
	response.SendResponse(w, resp)
}
