package tests

import (
	"fmt"
	"testing"

	test_utils "github.com/Bupher-Co/bupher-api/tests/utils"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	ts *test_utils.TestServer
}

func (s *AuthHandlerTestSuite) SetupTest() {
	s.ts = test_utils.NewTestServer()
}

func (s *AuthHandlerTestSuite) TearDownTest() {
	fmt.Println(s.ts)
	s.ts.DropTablesAndTypes()
	s.ts.Config.GetDB().Close()
}

func (s *AuthHandlerTestSuite) TestSignup() {}

// func (s *AuthHandlerTestSuite) TestSignin() {}

// func (s *AuthHandlerTestSuite) TestResetPassword() {}

func TestAuthHandlersSuite(t *testing.T) {
	suite.Run(t, &AuthHandlerTestSuite{})
}
