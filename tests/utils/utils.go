package utils

import (
	"context"
	"net/http/httptest"
	"time"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/routes"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/tests/utils/mocks/test_config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const createTypesSql = `
	CREATE TYPE ACCOUNT_TYPE_ENUM AS ENUM ('personal', 'business');
	CREATE TYPE EVENT_ENVIRONMENT_ENUM AS ENUM ('app_environment', 'push_environment', 'job_environment');
	CREATE TYPE EVENT_TYPE_ENUM AS ENUM ('sms', 'email');
	CREATE TYPE TOKEN_TYPE_ENUM AS ENUM ('access_token', 'refresh_token');

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
`

func createTypes(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, createTypesSql)
	return err
}

type testServer struct {
	Server *httptest.Server
	Config config.IConfig
}

func NewTestServer() *testServer {
	c := test_config.NewTestConfig()

	r := routes.GetRouter(c)

	s := httptest.NewServer(r)

	if err := createTypes(c.DB); err != nil {
		panic(err)
	}

	return &testServer{s, c}
}

func (ts *testServer) DropTables() {
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`

	ctx := context.Background()

	rows, err := ts.Config.GetDB().Query(ctx, query)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			panic(err)
		}

		_, err := ts.Config.GetDB().Exec(ctx, `DROP TABLE %s CASCADE`, tableName)
		if err != nil {
			panic(err)
		}
	}
}
