package tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/json"
	test_utils "github.com/Bupher-Co/bupher-api/tests/utils"
	"github.com/stretchr/testify/suite"
)

const (
	contentType = "application/json"
)

type TestUser struct {
	Email        string
	AccountType  string
	PhoneNumber  string
	BusinessName string
	FirstName    string
	LastName     string
	Password     string
	RegStage     int
}

type AuthHandlerTestSuite struct {
	suite.Suite
	ts       *test_utils.TestServer
	testUser TestUser
}

func (s *AuthHandlerTestSuite) SetupSuite() {
	s.ts = test_utils.NewTestServer()
	s.testUser = TestUser{
		Email:        "test1@user.com",
		AccountType:  "personal",
		PhoneNumber:  "09012345678",
		BusinessName: "Edeh Ventures",
		FirstName:    "Chimezie",
		LastName:     "Edeh",
		Password:     "password",
		RegStage:     1,
	}
}

func (s *AuthHandlerTestSuite) TearDownSuite() {
	s.ts.DropTablesAndTypes()
	s.ts.Config.GetDB().Close()
}

func (s *AuthHandlerTestSuite) TestSignupPersonal() {
	type DataResponse struct {
		Code string      `json:"code"`
		User models.User `json:"user"`
	}
	type MetaResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	type SuccessResponse struct {
		Success    bool         `json:"success"`
		Message    string       `json:"message"`
		StatusCode int          `json:"status_code"`
		Data       DataResponse `json:"data"`
		Meta       MetaResponse `json:"meta"`
	}

	// var verifyEmailCode string
	s.Run("Phase 1 sign up", func() {
		payload, _ := json.WriteJSON(map[string]any{
			"email":        s.testUser.Email,
			"account_type": s.testUser.AccountType,
			"reg_stage":    s.testUser.RegStage,
		})

		url := s.ts.Server.URL + "/api/v1/auth/sign-up"
		res, err := s.ts.Server.Client().Post(url, contentType, bytes.NewBuffer(payload))

		s.NoError(err)
		s.Equal(http.StatusOK, res.StatusCode)

		respBody := new(SuccessResponse)
		err = json.ReadJSON(res.Body, respBody)
		if err != nil {
			s.Fail(err.Error())
		}

		s.Equal(http.StatusOK, res.StatusCode)
		s.NotEmpty(respBody.Data.User)
		s.Equal(s.testUser.Email, respBody.Data.User.Email)
		s.NotEmpty(respBody.Data.Code)

		// verifyEmailCode = respBody.Data.Code
	})

	// phase 2 sign up
	// verify phone number

	// phase 3 sign up
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
