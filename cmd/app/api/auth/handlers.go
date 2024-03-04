package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/json"
	"github.com/Bupher-Co/bupher-api/pkg/jwt"
	"github.com/Bupher-Co/bupher-api/pkg/push"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/Bupher-Co/bupher-api/pkg/validator"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

const (
	VerifyEmailSubject = "Email verification"
	RegStage1Msg       = "verify your email"
	RegStage2Msg       = "verify your phone number"
	RegStage3Msg       = "sign up successful"
)

type authHandler struct {
	c *config.Config
}

func (h *authHandler) signUp(w http.ResponseWriter, r *http.Request) {
	body := new(signUpDto)
	resp := response.ApiResponse{}

	err := json.ReadJSON(r, body)
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

	user := new(models.User)
	user, err = h.c.UserRepository.GetByEmail(*body.Email, nil)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	if user != nil {
		switch user.RegStage {
		case int(utils.RegStage1):
			if !user.IsEmailVerified {
				otp := &models.Otp{
					UserID:    user.ID,
					Code:      utils.GenerateRandomNumber(),
					IsUsed:    false,
					OtpType:   models.EmailOtpType,
					ExpiresIn: time.Now().Add(models.OtpExpiresIn * time.Minute),
				}

				err = h.c.OtpRepository.Create(otp, nil)
				if err != nil {
					resp.Message = err.Error()
					response.SendErrorResponse(w, resp, http.StatusInternalServerError)
					return
				}

				utils.Background(func() {
					err = push.SendEmail(&push.Email{
						To:      []string{user.Email},
						Subject: VerifyEmailSubject,
						Text:    fmt.Sprintf("Use code %s to verify your email", otp.Code),
						Html:    fmt.Sprintf("<p>Use code %s to verify your email</p>", otp.Code),
					})

					if err != nil {
						h.c.Logger.Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
					}
				})

				resp.Message = RegStage1Msg
				if h.c.Env.IsDevelopment() {
					resp.Data = map[string]any{
						"code": otp.Code,
						"user": user,
					}
				} else {
					resp.Data = map[string]any{
						"user": user,
					}
				}

				response.SendResponse(w, resp)
				return
			}

		case int(utils.RegStage2):
			if !user.IsPhoneNumberVerified {
				otp := &models.Otp{
					UserID:    user.ID,
					Code:      utils.GenerateRandomNumber(),
					IsUsed:    false,
					OtpType:   models.SmsOtpType,
					ExpiresIn: time.Now().Add(models.OtpExpiresIn * time.Minute),
				}

				err = h.c.OtpRepository.Create(otp, nil)
				if err != nil {
					resp.Message = err.Error()
					response.SendErrorResponse(w, resp, http.StatusInternalServerError)
					return
				}

				utils.Background(func() {
					push.SendSMS(&push.Sms{
						Phone:   user.PhoneNumber.String,
						Message: fmt.Sprintf("Use code %s to verify your phone number", otp.Code),
					})
				})

				resp.Message = RegStage2Msg
				if h.c.Env.IsDevelopment() {
					resp.Data = map[string]string{
						"code": otp.Code,
					}
				}

				response.SendResponse(w, resp)
				return
			}

		case int(utils.RegStage3):
			resp.Message = "account already exists"
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}
	}

	tx, err := h.c.DB.Begin(context.Background())
	defer tx.Rollback(context.Background())

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	switch *body.RegStage {
	case utils.RegStage1:
		user = &models.User{
			Email:       *body.Email,
			AccountType: *body.AccountType,
			RegStage:    int(*body.RegStage),
		}

		err := h.c.UserRepository.Create(user, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		if *body.AccountType == "business" {
			business := &models.Business{
				UserID: user.ID,
				Name:   *body.BusinessName,
				Email:  *body.Email,
			}

			err = h.c.BusinessRepository.Create(business, tx)
			if err != nil {
				resp.Message = err.Error()
				response.SendErrorResponse(w, resp, http.StatusBadRequest)
				return
			}
		}

		otp := &models.Otp{
			UserID:    user.ID,
			Code:      utils.GenerateRandomNumber(),
			IsUsed:    false,
			OtpType:   models.EmailOtpType,
			ExpiresIn: time.Now().Add(models.OtpExpiresIn * time.Minute),
		}

		err = h.c.OtpRepository.Create(otp, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		utils.Background(func() {
			err = push.SendEmail(&push.Email{
				To:      []string{user.Email},
				Subject: VerifyEmailSubject,
				Text:    fmt.Sprintf("Use code %s to verify your email", otp.Code),
				Html:    fmt.Sprintf("<p>Use code %s to verify your email</p>", otp.Code),
			})

			if err != nil {
				h.c.Logger.Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
			}
		})

		resp.Message = RegStage1Msg

		if h.c.Env.IsDevelopment() {
			resp.Data = map[string]any{
				"code": otp.Code,
				"user": user,
			}
		} else {
			resp.Data = user
		}
	case utils.RegStage2:
		user.PhoneNumber = models.NullString{NullString: sql.NullString{String: *body.PhoneNumber}}
		user.RegStage = int(*body.RegStage)

		err = h.c.UserRepository.Update(user, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		otp := &models.Otp{
			UserID:    user.ID,
			Code:      utils.GenerateRandomNumber(),
			IsUsed:    false,
			OtpType:   models.SmsOtpType,
			ExpiresIn: time.Now().Add(models.OtpExpiresIn * time.Minute),
		}
		err = h.c.OtpRepository.Create(otp, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		utils.Background(func() {
			push.SendSMS(&push.Sms{
				Phone:   user.PhoneNumber.String,
				Message: fmt.Sprintf("Use code %s to verify your phone number", otp.Code),
			})
		})

		resp.Message = RegStage2Msg
		if h.c.Env.IsDevelopment() {
			resp.Data = map[string]any{
				"code": otp.Code,
				"user": user,
			}
		} else {
			resp.Data = user
		}
	case utils.RegStage3:
		user.FirstName = models.NullString{NullString: sql.NullString{String: *body.FirstName}}
		user.LastName = models.NullString{NullString: sql.NullString{String: *body.LastName}}
		user.RegStage = int(*body.RegStage)

		err = h.c.UserRepository.Update(user, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		hashPwd, err := utils.GenerateHash(*body.Password)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		auth := &models.Auth{
			UserID:   user.ID,
			Password: hashPwd,
		}

		err = h.c.AuthRepository.Create(auth, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		accessTokenStr, err := jwt.GenerateToken(&jwt.TokenClaims{
			UserID:    user.ID.String(),
			Email:     user.Email,
			TokenType: string(models.AccessToken),
		})
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		refreshTokenStr, err := jwt.GenerateToken(&jwt.TokenClaims{
			UserID:    user.ID.String(),
			Email:     user.Email,
			TokenType: string(models.RefreshToken),
		})
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		accessTokenHash, err := utils.GenerateHash(accessTokenStr)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		refreshTokenHash, err := utils.GenerateHash(accessTokenStr)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		accessToken := &models.Token{
			Hash:      accessTokenHash,
			UserID:    user.ID,
			InUse:     true,
			TokenType: models.AccessToken,
		}
		refreshToken := &models.Token{
			Hash:      refreshTokenHash,
			UserID:    user.ID,
			InUse:     true,
			TokenType: models.RefreshToken,
		}

		err = h.c.TokenRepository.Create(accessToken, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		err = h.c.TokenRepository.Create(refreshToken, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		resp.Message = RegStage3Msg
		resp.Data = user
		resp.Meta = response.ApiResponseMeta{
			AccessToken:  &accessTokenStr,
			RefreshToken: &refreshTokenStr,
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	response.SendResponse(w, resp)
}

func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *authHandler) verifyCode(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(verifyCodeDto)

	err := json.ReadJSON(r, body)
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

	user, err := h.c.UserRepository.GetByEmail(body.Email, nil)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = response.ErrNotFound.Error()
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	tx, err := h.c.DB.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	otp, err := h.c.OtpRepository.GetOneByQuery(`
		WHERE
			code = $1
			AND is_used = $2
			AND user_id = $3
			AND expires_in >= $4
	`, []any{body.Code, false, user.ID, time.Now()}, tx)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = "invalid or expired code"
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	otp.IsUsed = true
	switch body.OtpType {
	case "SMS":
		user.IsPhoneNumberVerified = true
		resp.Message = "phone number verified successfully"
	case "EMAIL":
		user.IsEmailVerified = true
		resp.Message = "email verified successfully"
	}

	err = h.c.UserRepository.Update(user, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = h.c.OtpRepository.Update(otp, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Data = user
	response.SendResponse(w, resp)
}

func (h *authHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *authHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}
