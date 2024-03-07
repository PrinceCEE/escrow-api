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
  hash BYTEA NOT NULL,
  token_type TOKEN_TYPE_ENUM NOT NULL,
  in_use BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ,
  version INT NOT NULL DEFAULT 1
);