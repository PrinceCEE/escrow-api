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

func (s *AuthHandlerTestSuite) TestSignup() {
	url := s.ts.Server.URL + "/api/v1/auth"
	post := s.ts.Server.Client().Post

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
			"otp_type": "EMAIL",
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
			"otp_type": "SMS",
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

		s.Equal(http.StatusOK, res.StatusCode)
		s.NotEmpty(respBody.Data.User)
		s.Equal(s.testBusiness.Email, respBody.Data.User.Email)
		s.NotEmpty(respBody.Data.Code)
	})
}

// func (s *AuthHandlerTestSuite) TestSignupBusiness() {
// 	// phase 1 sign up
// 	// verify email
// 	s.Run("Phase 1 sign up", func() {

// 	})

// 	// phase 2 sign up
// 	// verify phone number

// 	// phase 3 sign up
// }

// func (s *AuthHandlerTestSuite) TestSignin() {}

// func (s *AuthHandlerTestSuite) TestResetPassword() {}

func TestAuthHandlersSuite(t *testing.T) {
	suite.Run(t, &AuthHandlerTestSuite{})
}
