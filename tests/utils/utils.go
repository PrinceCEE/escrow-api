package test_utils

import (
	"context"
	"net/http/httptest"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/routes"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/tests/utils/mocks/test_config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const setupTypesSql = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TYPE ACCOUNT_TYPE_ENUM AS ENUM ('personal', 'business');
	CREATE TYPE EVENT_ENVIRONMENT_ENUM AS ENUM ('app_environment', 'push_environment', 'job_environment');
	CREATE TYPE EVENT_TYPE_ENUM AS ENUM ('sms', 'email');
	CREATE TYPE TOKEN_TYPE_ENUM AS ENUM ('access_token', 'refresh_token');
	CREATE TYPE OTP_TYPE AS ENUM ('SMS', 'EMAIL');

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		email VARCHAR(255) NOT NULL UNIQUE,
		phone_number VARCHAR(20) UNIQUE,
		first_name VARCHAR(80),
		last_name VARCHAR(80),
		is_phone_number_verified BOOLEAN DEFAULT false,
		is_email_verified BOOLEAN DEFAULT false,
		reg_stage INT CHECK (reg_stage IN (1, 2, 3)) NOT NULL,
		account_type ACCOUNT_TYPE_ENUM NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS businesses (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID REFERENCES users NOT NULL,
		name VARCHAR(80) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		deleted_at TIMESTAMPTZ,
		version INT NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS auths (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID REFERENCES users NOT NULL,
		password BYTEA NOT NULL,
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
`

var tearDownTypesSql = `
	DROP TABLE IF EXISTS tokens;
	DROP TABLE IF EXISTS events;
	DROP TABLE IF EXISTS auths;
	DROP TABLE IF EXISTS businesses;
	DROP TABLE IF EXISTS otps;
	DROP TABLE IF EXISTS users;

	DROP TYPE IF EXISTS ACCOUNT_TYPE_ENUM;
	DROP TYPE IF EXISTS EVENT_ENVIRONMENT_ENUM;
	DROP TYPE IF EXISTS EVENT_TYPE_ENUM;
	DROP TYPE IF EXISTS TOKEN_TYPE_ENUM;
	DROP TYPE IF EXISTS OTP_TYPE;
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

	if err := createTablesAndTypes(c.DB); err != nil {
		panic(err)
	}

	return &TestServer{s, c}
}

func (ts *TestServer) DropTablesAndTypes() {
	_, err := ts.Config.GetDB().Exec(context.Background(), tearDownTypesSql)
	if err != nil {
		panic(err)
	}
}
