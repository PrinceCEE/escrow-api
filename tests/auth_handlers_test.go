package tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/Bupher-Co/bupher-api/pkg/json"
	test_utils "github.com/Bupher-Co/bupher-api/tests/utils"
	"github.com/stretchr/testify/suite"
)

const (
	contentType = "application/json"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	ts           *test_utils.TestServer
	testUser     test_utils.TestUser
	testBusiness test_utils.TestBussiness
	password     string
}

type DataResponse struct {
	Code string              `json:"code"`
	User test_utils.TestUser `json:"user,omitempty"`
}

type MetaResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    DataResponse `json:"data,omitempty"`
	Meta    MetaResponse `json:"meta,omitempty"`
}

func (s *AuthHandlerTestSuite) SetupSuite() {
	s.ts = test_utils.NewTestServer()
	s.password = "password"
	s.testUser = test_utils.TestUser{
		Email:       "test1@user.com",
		AccountType: "personal",
		PhoneNumber: "09012345678",
		FirstName:   "Chimezie",
		LastName:    "Edeh",
		RegStage:    1,
	}
	s.testBusiness = test_utils.TestBussiness{Name: "Edeh Ventures", Email: "test1@business.com"}
}

func (s *AuthHandlerTestSuite) TearDownSuite() {
	s.ts.DropTablesAndTypes()
	s.ts.Config.GetDB().Close()
}

func (s *AuthHandlerTestSuite) TestAuthHandler() {
	url := s.ts.Server.URL + "/api/v1/auth"
	post := s.ts.Server.Client().Post

	// test sign up
	s.Run("sign up", func() {
		var verifyCode string

		s.Run("phase 1 sign up", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":        s.testUser.Email,
				"account_type": s.testUser.AccountType,
				"reg_stage":    s.testUser.RegStage,
			})

			res, err := post(url+"/sign-up", contentType, bytes.NewBuffer(payload))

			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()
			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.NotEmpty(respBody.Data.User)
			s.Equal(s.testUser.Email, respBody.Data.User.Email)
			s.NotEmpty(respBody.Data.Code)

			verifyCode = respBody.Data.Code
		})

		s.Run("should verify email", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":    s.testUser.Email,
				"code":     verifyCode,
				"otp_type": "email",
			})

			res, err := post(url+"/verify-code", contentType, bytes.NewBuffer(payload))

			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()
			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal("email verified successfully", respBody.Message)
		})

		s.Run("phase 2 sign up", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":        s.testUser.Email,
				"reg_stage":    2,
				"phone_number": s.testUser.PhoneNumber,
			})

			res, err := post(url+"/sign-up", contentType, bytes.NewBuffer(payload))

			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()
			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal(2, respBody.Data.User.RegStage)
			s.NotEmpty(respBody.Data.Code)

			verifyCode = respBody.Data.Code
		})

		s.Run("should verify phone", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":    s.testUser.Email,
				"code":     verifyCode,
				"otp_type": "sms",
			})

			res, err := post(url+"/verify-code", contentType, bytes.NewBuffer(payload))

			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()
			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal("phone number verified successfully", respBody.Message)
		})

		s.Run("phase 3 sign up", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":      s.testUser.Email,
				"first_name": s.testUser.FirstName,
				"last_name":  s.testUser.LastName,
				"password":   s.password,
				"reg_stage":  3,
			})

			res, err := post(url+"/sign-up", contentType, bytes.NewBuffer(payload))

			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()
			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.NotEmpty(respBody.Meta.AccessToken)
			s.NotEmpty(respBody.Meta.RefreshToken)
			s.Equal(3, respBody.Data.User.RegStage)
			s.Equal(s.testUser.FirstName, respBody.Data.User.FirstName)
		})

		s.Run("expects account exists error", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":        s.testUser.Email,
				"account_type": s.testUser.AccountType,
				"reg_stage":    s.testUser.RegStage,
			})

			res, err := post(url+"/sign-up", contentType, bytes.NewBuffer(payload))

			s.NoError(err)
			s.Equal(http.StatusBadRequest, res.StatusCode)

			defer res.Body.Close()
			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal(false, respBody.Success)
			s.Equal("account already exists", respBody.Message)
		})

		s.Run("phase 1 sign up for business", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":         s.testBusiness.Email,
				"account_type":  "business",
				"reg_stage":     1,
				"business_name": s.testBusiness.Name,
			})

			res, err := post(url+"/sign-up", contentType, bytes.NewBuffer(payload))
			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()
			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.NotEmpty(respBody.Data.User)
			s.Equal(s.testBusiness.Email, respBody.Data.User.Email)
			s.NotEmpty(respBody.Data.Code)
		})
	})

	// test Sign in
	s.Run("sign in", func() {
		s.Run("sign in personal account", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":    s.testUser.Email,
				"password": s.password,
			})

			res, err := post(url+"/sign-in", contentType, bytes.NewBuffer(payload))
			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()

			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.NotEmpty(respBody.Meta.AccessToken)
			s.NotEmpty(respBody.Meta.RefreshToken)
		})

		s.Run("sign in with wrong password", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":    s.testUser.Email,
				"password": "Pa55word",
			})

			res, err := post(url+"/sign-in", contentType, bytes.NewBuffer(payload))
			s.NoError(err)
			s.Equal(http.StatusBadRequest, res.StatusCode)

			defer res.Body.Close()

			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal("invalid sign in credentials", respBody.Message)
			s.Empty(respBody.Meta)
		})

		s.Run("sign in unverified email account", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":    s.testBusiness.Email,
				"password": s.password,
			})

			res, err := post(url+"/sign-in", contentType, bytes.NewBuffer(payload))
			s.NoError(err)
			s.Equal(http.StatusForbidden, res.StatusCode)

			defer res.Body.Close()

			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal("email not verified", respBody.Message)
			s.Empty(respBody.Meta)
		})
	})

	// test forgot and reset password
	s.Run("reset and forgot password", func() {
		var verifyCode string

		s.Run("forgot password", func() {
			payload, _ := json.WriteJSON(map[string]string{
				"email": s.testUser.Email,
			})

			res, err := post(url+"/reset-password", contentType, bytes.NewBuffer(payload))

			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()

			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal("otp send to your email", respBody.Message)
			s.NotEmpty(respBody.Data.Code)

			verifyCode = respBody.Data.Code
		})

		s.Run("should verify otp", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":    s.testUser.Email,
				"code":     verifyCode,
				"otp_type": "reset_password",
			})

			res, err := post(url+"/verify-code", contentType, bytes.NewBuffer(payload))

			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()
			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal("otp verified successfully", respBody.Message)
		})

		s.Run("change password with old one", func() {
			payload, _ := json.WriteJSON(map[string]string{
				"email":    s.testUser.Email,
				"password": s.password,
			})

			res, err := post(url+"/change-password", contentType, bytes.NewBuffer(payload))
			s.NoError(err)
			s.Equal(http.StatusBadRequest, res.StatusCode)

			defer res.Body.Close()

			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal("you can't use your old password", respBody.Message)
		})

		s.Run("change password", func() {
			payload, _ := json.WriteJSON(map[string]string{
				"email":    s.testUser.Email,
				"password": "password!",
			})

			res, err := post(url+"/change-password", contentType, bytes.NewBuffer(payload))
			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()

			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.Equal("password changed successfully", respBody.Message)
		})

		s.Run("sign in with new password", func() {
			payload, _ := json.WriteJSON(map[string]any{
				"email":    s.testUser.Email,
				"password": "password!",
			})

			res, err := post(url+"/sign-in", contentType, bytes.NewBuffer(payload))
			s.NoError(err)
			s.Equal(http.StatusOK, res.StatusCode)

			defer res.Body.Close()

			respBody := new(Response)
			_ = json.ReadJSON(res.Body, respBody)

			s.NotEmpty(respBody.Meta.AccessToken)
			s.NotEmpty(respBody.Meta.RefreshToken)
		})
	})
}

func TestAuthHandlersSuite(t *testing.T) {
	suite.Run(t, &AuthHandlerTestSuite{})
}
