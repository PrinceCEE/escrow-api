package auth

import (
	"net/http"

	appPkg "github.com/Bupher-Co/bupher-api/cmd/app/pkg"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/pkg/json"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/Bupher-Co/bupher-api/pkg/validator"
	"github.com/rs/zerolog"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	var body *signUpDto

	err := json.ReadJSON(r, body)
	if err != nil {
		config.Config.Logger.Log(zerolog.InfoLevel, "error reading request body", nil, err)
		appPkg.SendErrorResponse(
			w,
			appPkg.ApiResponse{Message: err.Error()},
			http.StatusBadRequest,
		)

		return
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		config.Config.Logger.Log(zerolog.InfoLevel, "signupDto validation error", validationErrors, appPkg.ErrBadRequest)
		appPkg.SendErrorResponse(
			w,
			appPkg.ApiResponse{
				Message: appPkg.ErrBadRequest.Error(),
				Data:    validationErrors,
			},
			http.StatusBadRequest,
		)

		return
	}

	// var user *models.User
	// check if user exists and has the same stage as in the input
	// if email isn't verified, send email
	// if phone number isn't verified, send sms
	// return the user
	switch *body.RegStage {
	case utils.RegStage1:
		// create user with email, account type and stage 1
		// if user is a business, then also create the business
		// then return the user
	case utils.RegStage2:
		// update the user phone number, and stage to 2
		// send sms
		// then return the user
	case utils.RegStage3:
		// update the user first name and last name and stage to 3
		// hash password and save the user's auth data
		// sign in the user (hash the token and save)
		// return both the user, access and refresh tokens
	}

	appPkg.SendErrorResponse(w, appPkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	appPkg.SendErrorResponse(w, appPkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func verifyCode(w http.ResponseWriter, r *http.Request) {
	appPkg.SendErrorResponse(w, appPkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	appPkg.SendErrorResponse(w, appPkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	appPkg.SendErrorResponse(w, appPkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}
