package test_utils

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/routes"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/pkg/json"
	"github.com/Bupher-Co/bupher-api/tests/utils/mocks/test_config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ContentType = "application/json"
)

const setupTypesSql = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TYPE ACCOUNT_TYPE_ENUM AS ENUM ('personal', 'business');
	CREATE TYPE EVENT_ENVIRONMENT_ENUM AS ENUM ('app_environment', 'push_environment', 'job_environment');
	CREATE TYPE EVENT_TYPE_ENUM AS ENUM ('sms', 'email');
	CREATE TYPE TOKEN_TYPE_ENUM AS ENUM ('access_token', 'refresh_token');
	CREATE TYPE OTP_TYPE AS ENUM ('sms', 'email', 'reset_password');
	CREATE TYPE MODEL_STATUS_ENUM AS ENUM (
		'Successful',
		'Canceled',
		'Pending'
	);
	CREATE TYPE WITHDRAWAL_TYPE_ENUM AS ENUM (
		'Withdrawal',
		'Deposit'
	);

	CREATE TABLE IF NOT EXISTS businesses (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(80) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		image_url TEXT,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		email VARCHAR(255) NOT NULL UNIQUE,
		phone_number VARCHAR(20) UNIQUE,
		first_name VARCHAR(80),
		last_name VARCHAR(80),
		image_url TEXT,
		is_phone_number_verified BOOLEAN DEFAULT false,
		is_email_verified BOOLEAN DEFAULT false,
		reg_stage INT CHECK (reg_stage IN (1, 2, 3)) NOT NULL,
		account_type ACCOUNT_TYPE_ENUM NOT NULL,
		business_id UUID REFERENCES businesses,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS auths (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID REFERENCES users NOT NULL UNIQUE,
		password TEXT NOT NULL,
		password_history JSON DEFAULT '[]',
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS events (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		data JSON,
		origin_environment EVENT_ENVIRONMENT_ENUM NOT NULL,
		target_environment EVENT_ENVIRONMENT_ENUM NOT NULL,
		event_type EVENT_TYPE_ENUM NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS tokens (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID REFERENCES users NOT NULL,
		hash TEXT NOT NULL,
		token_type TOKEN_TYPE_ENUM NOT NULL,
		in_use BOOLEAN DEFAULT true,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS otps (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users,
		code CHAR(4) NOT NULL,
		is_used BOOLEAN NOT NULL,
		otp_type OTP_TYPE NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		expires_in TIMESTAMPTZ NOT NULL,
		version INT DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS wallets (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		identifier UUID NOT NULL,
		balance INT NOT NULL DEFAULT 0,
		receivable_balance INT NOT NULL DEFAULT 0,
		payable_balance INT NOT NULL DEFAULT 0,
		account_type ACCOUNT_TYPE_ENUM NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS wallet_histories (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		wallet_id UUID REFERENCES wallets NOT NULL,
		type WITHDRAWAL_TYPE_ENUM NOT NULL,
		amount INT NOT NULL,
		status MODEL_STATUS_ENUM NOT NULL DEFAULT 'Pending',
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS bank_accounts (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		wallet_id UUID NOT NULL,
		bank_name VARCHAR(255) NOT NULL,
		account_name VARCHAR(255) NOT NULL,
		account_number CHAR(10) NOT NULL,
		bvn CHAR(11) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT DEFAULT 1,

		CONSTRAINT unique_bankname_accountnumber UNIQUE(account_number,bank_name)
	);
`

var tearDownTypesSql = `
	DROP TABLE IF EXISTS tokens;
	DROP TABLE IF EXISTS events;
	DROP TABLE IF EXISTS auths;
	DROP TABLE IF EXISTS otps;
	DROP TABLE IF EXISTS wallet_histories;
	DROP TABLE IF EXISTS bank_accounts;
	DROP TABLE IF EXISTS wallets;
	DROP TABLE IF EXISTS users;
	DROP TABLE IF EXISTS businesses;

	DROP TYPE IF EXISTS ACCOUNT_TYPE_ENUM;
	DROP TYPE IF EXISTS EVENT_ENVIRONMENT_ENUM;
	DROP TYPE IF EXISTS EVENT_TYPE_ENUM;
	DROP TYPE IF EXISTS TOKEN_TYPE_ENUM;
	DROP TYPE IF EXISTS OTP_TYPE;
	DROP TYPE IF EXISTS MODEL_STATUS_ENUM;
	DROP TYPE IF EXISTS WITHDRAWAL_TYPE_ENUM;
`

func createTablesAndTypes(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), setupTypesSql)
	return err
}

type TestServer struct {
	Server *httptest.Server
	Config config.IConfig
}

func NewTestServer() *TestServer {
	c := test_config.NewTestConfig()

	r := routes.GetRouter(c)

	s := httptest.NewServer(r)

	ts := TestServer{s, c}
	ts.DropTablesAndTypes()

	if err := createTablesAndTypes(c.DB); err != nil {
		panic(err)
	}

	return &ts
}

func (ts *TestServer) DropTablesAndTypes() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := ts.Config.GetDB().Exec(ctx, tearDownTypesSql)

	if err != nil {
		panic(err)
	}
}

type TestBankAccount struct {
	BankName      string     `json:"bank_name"`
	AccountName   string     `json:"account_name"`
	AccountNumber string     `json:"account_number"`
	BVN           string     `json:"bvn" db:"bvn"`
	WalletID      string     `json:"wallet_id"`
	Wallet        TestWallet `json:"wallet,omitempty"`
	TestModelMixin
}

type TestWallet struct {
	Balance     int           `json:"balance"`
	Receivable  int           `json:"receivable_balance"`
	Payable     int           `json:"payable_balance"`
	AccountType string        `json:"account_type"`
	Identifier  string        `json:"identifier"`
	User        *TestUser     `json:"user,omitempty"`
	Business    *TestBusiness `json:"business,omitempty"`
	TestModelMixin
}

type TestWalletHistory struct {
	WalletID string      `json:"wallet_id"`
	Type     string      `json:"type"`
	Amount   int         `json:"amount"`
	Status   string      `json:"status"`
	Wallet   *TestWallet `json:"wallet,omitempty"`
	TestModelMixin
}

type MetaResponse struct {
	Page         int    `json:"page,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	Total        int    `json:"total,omitempty"`
	TotalPages   int    `json:"total_pages,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type SignupDataResponse struct {
	Code string   `json:"code"`
	User TestUser `json:"user,omitempty"`
}

type GetUserResponse struct {
	User TestUser `json:"user,omitempty"`
}

type Response[T any] struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    T            `json:"data,omitempty"`
	Meta    MetaResponse `json:"meta,omitempty"`
}

type TestModelMixin struct {
	ID        string     `json:"id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version,omitempty"`
}

type TestUser struct {
	TestModelMixin
	Email                 string        `json:"email,omitempty"`
	PhoneNumber           string        `json:"phone_number,omitempty"`
	FirstName             string        `json:"first_name,omitempty"`
	LastName              string        `json:"last_name,omitempty"`
	IsPhoneNumberVerified bool          `json:"is_phone_number_verified,omitempty"`
	IsEmailVerified       bool          `json:"is_email_verified,omitempty"`
	RegStage              int           `json:"reg_stage,omitempty"`
	AccountType           string        `json:"account_type,omitempty"`
	ImageUrl              string        `json:"image_url,omitempty"`
	BusinessID            *string       `json:"business_id,omitempty"`
	Business              *TestBusiness `json:"business,omitempty"`
}

type TestBusiness struct {
	UserID   string `json:"user_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	ImageUrl string `json:"image_url,omitempty"`
	TestModelMixin
}

func SignupPersonalUser(ts *TestServer) (TestUser, string) {
	url := ts.Server.URL + "/api/v1/auth"
	post := ts.Server.Client().Post
	contentType := "application/json"
	email := "testuser2@user.com"

	// phase 1 sign up
	phase1SignupDto := map[string]any{
		"email":        email,
		"reg_stage":    1,
		"account_type": "personal",
	}

	data, _ := json.Marshal(phase1SignupDto)
	res, err := post(url+"/sign-up", contentType, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	respBody := new(Response[SignupDataResponse])
	_ = json.ReadJSON(res.Body, respBody)
	res.Body.Close()

	if !respBody.Success {
		panic(respBody.Message)
	}

	// email verification
	verifyCodeDto := map[string]any{
		"email":    email,
		"code":     respBody.Data.Code,
		"otp_type": "email",
	}

	data, _ = json.Marshal(verifyCodeDto)
	res, err = post(url+"/verify-code", contentType, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	respBody = new(Response[SignupDataResponse])
	_ = json.ReadJSON(res.Body, respBody)
	res.Body.Close()

	if !respBody.Success {
		panic(respBody.Message)
	}

	// phase 2 sign up
	phase2Signup := map[string]any{
		"email":        email,
		"phone_number": "09012345678",
		"reg_stage":    2,
	}

	data, _ = json.Marshal(phase2Signup)
	res, err = http.Post(url+"/sign-up", contentType, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	respBody = new(Response[SignupDataResponse])
	_ = json.ReadJSON(res.Body, respBody)
	res.Body.Close()

	if !respBody.Success {
		panic(respBody.Message)
	}

	// phone number verification
	verifyCodeDto = map[string]any{
		"email":    email,
		"code":     respBody.Data.Code,
		"otp_type": "sms",
	}

	data, _ = json.Marshal(verifyCodeDto)
	res, err = post(url+"/verify-code", contentType, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	respBody = new(Response[SignupDataResponse])
	_ = json.ReadJSON(res.Body, respBody)
	res.Body.Close()

	if !respBody.Success {
		panic(respBody.Message)
	}

	// phase 3 sign up
	phase3Signup := map[string]any{
		"email":      email,
		"first_name": "Test",
		"last_name":  "User",
		"password":   "password",
		"reg_stage":  3,
	}

	data, _ = json.Marshal(phase3Signup)
	res, err = http.Post(url+"/sign-up", contentType, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	respBody = new(Response[SignupDataResponse])
	_ = json.ReadJSON(res.Body, respBody)
	res.Body.Close()

	if !respBody.Success {
		panic(respBody.Message)
	}

	return respBody.Data.User, respBody.Meta.AccessToken
}
