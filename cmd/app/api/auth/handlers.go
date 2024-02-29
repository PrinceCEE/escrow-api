package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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

func signUp(w http.ResponseWriter, r *http.Request) {
	var body *signUpDto
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

	var user *models.User
	user, err = config.Config.Repositories.UserRepository.GetByEmail(*body.Email, nil)
	if err != nil {
		var statusCode int

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusBadRequest
		}

		config.Config.Logger.Log(zerolog.InfoLevel, "error fetching user", nil, err)
		resp.Message = response.ErrNotFound.Error()
		response.SendErrorResponse(w, resp, statusCode)
		return
	}

	if user != nil {
		switch user.RegStage {
		case int(utils.RegStage1):
			if !user.IsEmailVerified {
				otp := &models.Otp{
					UserID:  user.ID,
					Code:    utils.GenerateRandomNumber(),
					IsUsed:  false,
					OtpType: models.EmailOtpType,
				}

				err = config.Config.Repositories.OtpRepository.Create(otp, nil)
				if err != nil {
					resp.Message = err.Error()
					response.SendErrorResponse(w, resp, http.StatusInternalServerError)
					return
				}

				push.SendEmail(&push.Email{
					To:      []string{*user.Email},
					Subject: VerifyEmailSubject,
					Text:    fmt.Sprintf("Use code %s to verify your email", otp.Code),
					Html:    fmt.Sprintf("<p>Use code %s to verify your email</p>", otp.Code),
				})

				resp.Message = RegStage1Msg
				if config.Config.Env.IsDevelopment() {
					resp.Data = map[string]string{
						"code": otp.Code,
					}
				}

				response.SendResponse(w, resp)
				return
			}

		case int(utils.RegStage2):
			if !user.IsPhoneNumberVerified {
				otp := &models.Otp{
					UserID:  user.ID,
					Code:    utils.GenerateRandomNumber(),
					IsUsed:  false,
					OtpType: models.SmsOtpType,
				}

				err = config.Config.Repositories.OtpRepository.Create(otp, nil)
				if err != nil {
					resp.Message = err.Error()
					response.SendErrorResponse(w, resp, http.StatusInternalServerError)
					return
				}

				push.SendSMS(&push.Sms{
					Phone:   *user.PhoneNumber,
					Message: fmt.Sprintf("Use code %s to verify your phone number", otp.Code),
				})

				resp.Message = RegStage2Msg
				if config.Config.Env.IsDevelopment() {
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

	tx, err := config.Config.DB.Begin(context.Background())
	defer tx.Rollback(context.Background())

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	switch *body.RegStage {
	case utils.RegStage1:
		user = &models.User{
			Email:       body.Email,
			AccountType: *body.AccountType,
			RegStage:    int(*body.RegStage),
		}

		err := config.Config.Repositories.UserRepository.Create(user, tx)
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

			err = config.Config.Repositories.BusinessRepository.Create(business, tx)
			if err != nil {
				resp.Message = err.Error()
				response.SendErrorResponse(w, resp, http.StatusBadRequest)
				return
			}
		}

		otp := &models.Otp{
			UserID:  user.ID,
			Code:    utils.GenerateRandomNumber(),
			IsUsed:  false,
			OtpType: models.EmailOtpType,
		}
		err = config.Config.Repositories.OtpRepository.Create(otp, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		push.SendEmail(&push.Email{
			To:      []string{*user.Email},
			Subject: VerifyEmailSubject,
			Text:    fmt.Sprintf("Use code %s to verify your email", otp.Code),
			Html:    fmt.Sprintf("<p>Use code %s to verify your email</p>", otp.Code),
		})

		resp.Message = RegStage1Msg

		if config.Config.Env.IsDevelopment() {
			resp.Data = map[string]any{
				"code": otp.Code,
				"user": user,
			}
		} else {
			resp.Data = user
		}
	case utils.RegStage2:
		user.PhoneNumber = body.PhoneNumber
		user.RegStage = int(*body.RegStage)

		err = config.Config.Repositories.UserRepository.Update(user, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		otp := &models.Otp{
			UserID:  user.ID,
			Code:    utils.GenerateRandomNumber(),
			IsUsed:  false,
			OtpType: models.SmsOtpType,
		}
		err = config.Config.Repositories.OtpRepository.Create(otp, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		push.SendSMS(&push.Sms{
			Phone:   *user.PhoneNumber,
			Message: fmt.Sprintf("Use code %s to verify your phone number", otp.Code),
		})

		resp.Message = RegStage2Msg
		if config.Config.Env.IsDevelopment() {
			resp.Data = map[string]any{
				"code": otp.Code,
				"user": user,
			}
		} else {
			resp.Data = user
		}
	case utils.RegStage3:
		user.FirstName = *body.FirstName
		user.LastName = *body.LastName
		user.RegStage = int(*body.RegStage)

		err = config.Config.Repositories.UserRepository.Update(user, tx)
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

		err = config.Config.Repositories.AuthRepository.Create(auth, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		accessTokenStr, err := jwt.GenerateToken(&jwt.TokenClaims{
			UserID:    user.ID,
			Email:     *user.Email,
			TokenType: string(models.AccessToken),
		})
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		refreshTokenStr, err := jwt.GenerateToken(&jwt.TokenClaims{
			UserID:    user.ID,
			Email:     *user.Email,
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

		err = config.Config.Repositories.TokenRepository.Create(accessToken, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		err = config.Config.Repositories.TokenRepository.Create(refreshToken, tx)
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

func signIn(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func verifyCode(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}
