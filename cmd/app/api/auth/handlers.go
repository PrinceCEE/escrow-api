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
	"golang.org/x/crypto/bcrypt"
)

const (
	VerifyEmailSubject = "Email verification"
	RegStage1Msg       = "verify your email"
	RegStage2Msg       = "verify your phone number"
	RegStage3Msg       = "sign up successful"
)

type IConfig interface{}

type authHandler struct {
	c config.IConfig
}

func (h *authHandler) signUp(w http.ResponseWriter, r *http.Request) {
	body := new(signUpDto)
	resp := response.ApiResponse{}
	env := h.c.Getenv("ENVIRONMENT")

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

	userRepo := h.c.GetUserRepository()
	businessRepo := h.c.GetBusinessRepository()
	authRepo := h.c.GetAuthRepository()
	otpRepo := h.c.GetOtpRepository()
	tokenRepo := h.c.GetTokenRepository()
	walletRepo := h.c.GetWalletRepository()

	user := new(models.User)
	user, err = userRepo.GetByEmail(*body.Email, nil)
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

				err = otpRepo.Create(otp, nil)
				if err != nil {
					resp.Message = err.Error()
					response.SendErrorResponse(w, resp, http.StatusInternalServerError)
					return
				}

				utils.Background(func() {
					err = h.c.GetPush().SendEmail(&push.Email{
						To:      []string{user.Email},
						Subject: VerifyEmailSubject,
						Text:    fmt.Sprintf("Use code %s to verify your email", otp.Code),
						Html:    fmt.Sprintf("<p>Use code %s to verify your email</p>", otp.Code),
					})

					if err != nil {
						h.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
					}
				})

				resp.Message = RegStage1Msg
				if env == "development" || env == "test" {
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

				err = otpRepo.Create(otp, nil)
				if err != nil {
					resp.Message = err.Error()
					response.SendErrorResponse(w, resp, http.StatusInternalServerError)
					return
				}

				utils.Background(func() {
					h.c.GetPush().SendSMS(&push.Sms{
						Phone:   user.PhoneNumber.String,
						Message: fmt.Sprintf("Use code %s to verify your phone number", otp.Code),
					})
				})

				resp.Message = RegStage2Msg
				if env == "development" || env == "test" {
					utils.Background(func() {
						err = h.c.GetPush().SendEmail(&push.Email{
							To:      []string{user.Email},
							Subject: VerifyEmailSubject,
							Text:    fmt.Sprintf("Use code %s to verify your phone number", otp.Code),
							Html:    fmt.Sprintf("<p>Use code %s to verify your phone number</p>", otp.Code),
						})

						if err != nil {
							h.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
						}
					})

					resp.Data = map[string]any{
						"code": otp.Code,
						"user": user,
					}
				} else {
					resp.Data = map[string]any{
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

	tx, err := h.c.GetDB().Begin(context.Background())
	defer tx.Rollback(context.Background())

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	switch *body.RegStage {
	case utils.RegStage1:
		var business *models.Business
		if *body.AccountType == "business" {
			business = &models.Business{
				Name:  *body.BusinessName,
				Email: *body.Email,
			}

			err = businessRepo.Create(business, tx)
			if err != nil {
				resp.Message = err.Error()
				response.SendErrorResponse(w, resp, http.StatusBadRequest)
				return
			}
		}

		user = &models.User{
			Email:       *body.Email,
			AccountType: *body.AccountType,
			RegStage:    int(*body.RegStage),
		}
		if business != nil {
			user.BusinessID = business.ID
			user.Business = business
		}

		err := userRepo.Create(user, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		otp := &models.Otp{
			UserID:    user.ID,
			Code:      utils.GenerateRandomNumber(),
			IsUsed:    false,
			OtpType:   models.EmailOtpType,
			ExpiresIn: time.Now().Add(models.OtpExpiresIn * time.Minute),
		}

		err = otpRepo.Create(otp, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		utils.Background(func() {
			err = h.c.GetPush().SendEmail(&push.Email{
				To:      []string{user.Email},
				Subject: VerifyEmailSubject,
				Text:    fmt.Sprintf("Use code %s to verify your email", otp.Code),
				Html:    fmt.Sprintf("<p>Use code %s to verify your email</p>", otp.Code),
			})

			if err != nil {
				h.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
			}
		})

		resp.Message = RegStage1Msg

		if env == "development" || env == "test" {
			resp.Data = map[string]any{
				"code": otp.Code,
				"user": user,
			}
		} else {
			resp.Data = map[string]any{
				"user": user,
			}
		}
	case utils.RegStage2:
		user.PhoneNumber = models.NullString{NullString: sql.NullString{String: *body.PhoneNumber, Valid: true}}
		user.RegStage = int(*body.RegStage)

		err = userRepo.Update(user, tx)
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
		err = otpRepo.Create(otp, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		utils.Background(func() {
			h.c.GetPush().SendSMS(&push.Sms{
				Phone:   user.PhoneNumber.String,
				Message: fmt.Sprintf("Use code %s to verify your phone number", otp.Code),
			})
		})

		resp.Message = RegStage2Msg
		if env == "development" || env == "test" {
			utils.Background(func() {
				err = h.c.GetPush().SendEmail(&push.Email{
					To:      []string{user.Email},
					Subject: VerifyEmailSubject,
					Text:    fmt.Sprintf("Use code %s to verify your phone number", otp.Code),
					Html:    fmt.Sprintf("<p>Use code %s to verify your phone number</p>", otp.Code),
				})

				if err != nil {
					h.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
				}
			})

			resp.Data = map[string]any{
				"code": otp.Code,
				"user": user,
			}
		} else {
			resp.Data = map[string]any{
				"user": user,
			}
		}
	case utils.RegStage3:
		user.FirstName = models.NullString{NullString: sql.NullString{String: *body.FirstName, Valid: true}}
		user.LastName = models.NullString{NullString: sql.NullString{String: *body.LastName, Valid: true}}
		user.RegStage = int(*body.RegStage)

		err = userRepo.Update(user, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		hashPwd, err := utils.GeneratePasswordHash(*body.Password)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		auth := &models.Auth{
			UserID:   user.ID,
			Password: string(hashPwd),
		}

		err = authRepo.Create(auth, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		wallet := &models.Wallet{AccountType: user.AccountType}
		if user.AccountType == models.PersonalAccountType {
			wallet.Identifier = user.ID
		} else {
			wallet.Identifier = user.BusinessID
		}

		err = walletRepo.Create(wallet, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		accessTokenStr, _ := jwt.GenerateToken(&jwt.TokenClaims{
			UserID:    user.ID,
			Email:     user.Email,
			TokenType: string(models.AccessToken),
		})

		refreshTokenStr, _ := jwt.GenerateToken(&jwt.TokenClaims{
			UserID:    user.ID,
			Email:     user.Email,
			TokenType: string(models.RefreshToken),
		})

		accessToken := &models.Token{
			Hash:      accessTokenStr,
			UserID:    user.ID,
			InUse:     true,
			TokenType: models.AccessToken,
		}
		refreshToken := &models.Token{
			Hash:      refreshTokenStr,
			UserID:    user.ID,
			InUse:     true,
			TokenType: models.RefreshToken,
		}

		err = tokenRepo.Create(accessToken, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		err = tokenRepo.Create(refreshToken, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		resp.Message = RegStage3Msg
		resp.Data = map[string]any{"user": user}
		resp.Meta = response.ApiResponseMeta{
			AccessToken:  accessTokenStr,
			RefreshToken: refreshTokenStr,
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
	resp := response.ApiResponse{}
	body := new(signInDto)

	userRepo := h.c.GetUserRepository()
	authRepo := h.c.GetAuthRepository()
	tokenRepo := h.c.GetTokenRepository()

	err := json.ReadJSON(r.Body, body)
	defer r.Body.Close()

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	user, err := userRepo.GetByEmail(body.Email, nil)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = "account doesn't exist"
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	if !user.IsEmailVerified {
		resp.Message = "email not verified"
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}
	if !user.IsPhoneNumberVerified {
		resp.Message = "phone number not verified"
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	auth, err := authRepo.GetByUserId(user.ID, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = utils.ComparePassword(body.Password, []byte(auth.Password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			resp.Message = "invalid sign in credentials"
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	accessTokenStr, _ := jwt.GenerateToken(&jwt.TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		TokenType: string(models.AccessToken),
	})

	refreshTokenStr, _ := jwt.GenerateToken(&jwt.TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		TokenType: string(models.RefreshToken),
	})

	accessToken := &models.Token{
		Hash:      accessTokenStr,
		UserID:    user.ID,
		InUse:     true,
		TokenType: models.AccessToken,
	}
	refreshToken := &models.Token{
		Hash:      refreshTokenStr,
		UserID:    user.ID,
		InUse:     true,
		TokenType: models.RefreshToken,
	}

	err = tokenRepo.Create(accessToken, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = tokenRepo.Create(refreshToken, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Message = "signed in successfully"
	resp.Meta = response.ApiResponseMeta{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}

	response.SendResponse(w, resp)
}

func (h *authHandler) resendCode(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(resendCodeOTPDto)
	env := h.c.Getenv("ENVIRONMENT")

	err := json.ReadJSON(r.Body, body)
	defer r.Body.Close()

	userRepo := h.c.GetUserRepository()
	otpRepo := h.c.GetOtpRepository()

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	user := new(models.User)
	if body.OtpType == models.SmsOtpType {
		user, err = userRepo.GetByPhoneNumber(body.Identifier, nil)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}
	} else {
		user, err = userRepo.GetByEmail(body.Identifier, nil)
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
		OtpType:   body.OtpType,
		ExpiresIn: time.Now().Add(models.OtpExpiresIn * time.Minute),
	}
	err = otpRepo.Create(otp, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	switch body.OtpType {
	case models.SmsOtpType:
		utils.Background(func() {
			h.c.GetPush().SendSMS(&push.Sms{
				Phone:   user.PhoneNumber.String,
				Message: fmt.Sprintf("Use code %s to verify your phone number", otp.Code),
			})
		})
	case models.EmailOtpType:
		utils.Background(func() {
			err = h.c.GetPush().SendEmail(&push.Email{
				To:      []string{user.Email},
				Subject: VerifyEmailSubject,
				Text:    fmt.Sprintf("Use code %s to verify your email", otp.Code),
				Html:    fmt.Sprintf("<p>Use code %s to verify your email</p>", otp.Code),
			})

			if err != nil {
				h.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
			}
		})
	case models.ResetPasswordType:
		utils.Background(func() {
			err = h.c.GetPush().SendEmail(&push.Email{
				To:      []string{user.Email},
				Subject: VerifyEmailSubject,
				Text:    fmt.Sprintf("Use code %s to verify your email", otp.Code),
				Html:    fmt.Sprintf("<p>Use code %s to verify your email</p>", otp.Code),
			})

			if err != nil {
				h.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
			}
		})
	}

	if body.OtpType == models.SmsOtpType {
		resp.Message = "OTP sent to your phone number"
	} else {
		resp.Message = "OTP sent to your email"
	}

	if env == "development" || env == "test" {
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
}

func (h *authHandler) verifyCode(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(verifyCodeDto)

	err := json.ReadJSON(r.Body, body)
	defer r.Body.Close()

	userRepo := h.c.GetUserRepository()
	otpRepo := h.c.GetOtpRepository()

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

	user, err := userRepo.GetByEmail(body.Email, nil)
	if err != nil {
		fmt.Println("Line 662", err)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = fmt.Sprintf("user %s", response.ErrNotFound.Error())
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	tx, err := h.c.GetDB().Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	otp, err := otpRepo.GetOneByWhere(`
		WHERE
			code = $1
			AND is_used = $2
			AND user_id = $3
			AND expires_in >= $4
			AND otp_type = $5
	`, []any{body.Code, false, user.ID, time.Now(), body.OtpType}, tx)

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
	case "sms":
		user.IsPhoneNumberVerified = true
		resp.Message = "phone number verified successfully"
	case "email":
		user.IsEmailVerified = true
		resp.Message = "email verified successfully"
	case "reset_password":
		resp.Message = "otp verified successfully"
	}

	err = userRepo.Update(user, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = otpRepo.Update(otp, tx)
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
	resp := response.ApiResponse{}
	body := new(forgotPasswordDto)
	defer r.Body.Close()

	err := json.ReadJSON(r.Body, body)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	userRepo := h.c.GetUserRepository()
	otpRepo := h.c.GetOtpRepository()

	user, err := userRepo.GetByEmail(body.Email, nil)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = "user not found"
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	if !user.IsEmailVerified {
		resp.Message = "email not verified"
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	code := utils.GenerateRandomNumber()
	otp := &models.Otp{
		UserID:    user.ID,
		Code:      code,
		IsUsed:    false,
		ExpiresIn: time.Now().Add(models.OtpExpiresIn * time.Minute),
		OtpType:   models.ResetPasswordType,
	}

	utils.Background(func() {
		err = h.c.GetPush().SendEmail(&push.Email{
			To:      []string{user.Email},
			Subject: VerifyEmailSubject,
			Text:    fmt.Sprintf("Use code %s to verify your email", otp.Code),
			Html:    fmt.Sprintf("<p>Use code %s to verify your email</p>", otp.Code),
		})

		if err != nil {
			h.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
		}
	})

	err = otpRepo.Create(otp, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Message = "otp send to your email"
	env := h.c.Getenv("ENVIRONMENT")
	if env == "development" || env == "test" {
		resp.Data = map[string]string{
			"code": code,
		}
	}

	response.SendResponse(w, resp)
}

func (h *authHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(changePasswordDto)

	err := json.ReadJSON(r.Body, body)
	defer r.Body.Close()

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	userRepo := h.c.GetUserRepository()
	authRepo := h.c.GetAuthRepository()

	user, err := userRepo.GetByEmail(body.Email, nil)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = "user not found"
		default:
			resp.Message = err.Error()
		}

		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	auth, err := authRepo.GetByUserId(user.ID, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = utils.ComparePassword(body.Password, []byte(auth.Password))
	if err == nil {
		resp.Message = "you can't use your old password"
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	pwdHash, err := utils.GeneratePasswordHash(body.Password)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	oldHash := auth.Password
	auth.Password = string(pwdHash)
	auth.PasswordHistory = append(auth.PasswordHistory, models.PasswordHistory{
		Password:  oldHash,
		Timestamp: time.Now(),
	})

	err = authRepo.Update(auth, nil)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Message = "password changed successfully"
	response.SendResponse(w, resp)
}
