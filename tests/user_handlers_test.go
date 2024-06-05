package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/princecee/escrow-api/pkg/json"
	test_utils "github.com/princecee/escrow-api/tests/utils"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	ts          *test_utils.TestServer
	user        test_utils.TestUser
	accessToken string
}

func (s *UserHandlerTestSuite) SetupSuite() {
	s.ts = test_utils.NewTestServer()

	user, token := test_utils.SignupPersonalUser(s.ts)
	s.user = user
	s.accessToken = token
}

func (s *UserHandlerTestSuite) TearDownSuite() {
	s.ts.DropTablesAndTypes()
	s.ts.Config.GetDB().Close()
}

func (s *UserHandlerTestSuite) get(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", s.accessToken)},
		"Content-Type":  {test_utils.ContentType},
	}

	return req
}

func (s *UserHandlerTestSuite) put(url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(http.MethodPut, url, body)
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", s.accessToken)},
		"Content-Type":  {test_utils.ContentType},
	}

	return req
}

func (s *UserHandlerTestSuite) TestUserHandler() {
	client := s.ts.Server.Client()
	url := s.ts.Server.URL + "/api/v1/users"

	var userId string

	s.Run("get me", func() {
		req := s.get(url + "/me")
		res, err := client.Do(req)
		s.NoError(err)

		respBody := new(test_utils.Response[test_utils.GetUserResponse])
		_ = json.ReadJSON(res.Body, respBody)
		defer res.Body.Close()

		s.Equal(true, respBody.Success)
		s.NotEmpty(respBody.Data.User)
		s.Equal(s.user.Email, respBody.Data.User.Email)

		userId = respBody.Data.User.ID
	})

	s.Run("get user by id", func() {
		req := s.get(url + "/" + userId)
		res, err := client.Do(req)
		s.NoError(err)

		respBody := new(test_utils.Response[test_utils.GetUserResponse])
		_ = json.ReadJSON(res.Body, respBody)
		defer res.Body.Close()

		s.Equal(true, respBody.Success)
		s.NotEmpty(respBody.Data.User)
		s.Equal(s.user.Email, respBody.Data.User.Email)
	})

	s.Run("update user", func() {
		updateAccountDto := map[string]string{
			"first_name": "TestTest",
			"last_name":  "UserUser",
			"image_url":  "https://www.google.com",
		}

		data, _ := json.Marshal(updateAccountDto)
		req := s.put(url+"/update-account", bytes.NewBuffer(data))

		res, err := client.Do(req)
		s.NoError(err)

		respBody := new(test_utils.Response[test_utils.GetUserResponse])
		_ = json.ReadJSON(res.Body, respBody)
		defer res.Body.Close()

		s.Equal(true, respBody.Success)
		s.Equal(http.StatusOK, res.StatusCode)
		s.NotEqual("User", respBody.Data.User.LastName)
		s.Equal(updateAccountDto["image_url"], respBody.Data.User.ImageUrl)
	})

	s.Run("change password", func() {
		password := "passwordsss"
		changePassword := map[string]string{"password": password}

		data, _ := json.Marshal(changePassword)
		req := s.put(url+"/change-password", bytes.NewBuffer(data))

		res, err := client.Do(req)
		s.NoError(err)

		respBody := new(test_utils.Response[any])
		_ = json.ReadJSON(res.Body, respBody)
		defer res.Body.Close()

		s.Equal(true, respBody.Success)
		s.Equal("password changed successfully", respBody.Message)

		s.Run("sign in with the changed password", func() {
			url := s.ts.Server.URL + "/api/v1/auth/sign-in"

			signInDto := map[string]string{
				"email":    s.user.Email,
				"password": password,
			}

			data, _ := json.Marshal(signInDto)
			res, err := s.ts.Server.Client().Post(url, test_utils.ContentType, bytes.NewBuffer(data))
			s.NoError(err)

			respBody := new(test_utils.Response[any])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(true, respBody.Success)
			s.Equal("signed in successfully", respBody.Message)
			s.NotEmpty(respBody.Meta.AccessToken)
			s.NotEmpty(respBody.Meta.RefreshToken)
		})

		s.Run("sign in with the old password", func() {
			url := s.ts.Server.URL + "/api/v1/auth/sign-in"

			signInDto := map[string]string{
				"email":    s.user.Email,
				"password": "password",
			}

			data, _ := json.Marshal(signInDto)
			res, err := s.ts.Server.Client().Post(url, test_utils.ContentType, bytes.NewBuffer(data))
			s.NoError(err)

			respBody := new(test_utils.Response[any])
			_ = json.ReadJSON(res.Body, respBody)
			defer res.Body.Close()

			s.Equal(false, respBody.Success)
			s.Equal("invalid sign in credentials", respBody.Message)
			s.Empty(respBody.Meta.AccessToken)
			s.Empty(respBody.Meta.RefreshToken)
		})
	})
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, &UserHandlerTestSuite{})
}
