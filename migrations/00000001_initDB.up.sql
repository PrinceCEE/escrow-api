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
CREATE TYPE TRANSACTION_TYPE_ENUM AS ENUM (
	'Product',
	'Service',
	'Crypto'
);
CREATE TYPE TRANSACTION_CREATED_BY_ENUM AS ENUM (
	'Seller',
	'Buyer'
);
CREATE TYPE TRANSACTION_STATUS_ENUM AS ENUM (
	'Sent-Awaiting',
	'Pending-Payment',
	'Pending-Delivery',
	'Canceled',
	'Completed'
);
CREATE TYPE TRANSACTION_TIMELINE_NAME_ENUM AS ENUM (
	'Transaction Created',
	'Transaction Approved',
	'Payment Submitted',
	'Delivery Done',
	'Marked As Completed',
	'Transaction Canceled'
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
  status MODEL_STATUS_ENUM NOT NULL DEFAULT "Pending",
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

CREATE TABLE IF NOT EXISTS transactions (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	status TRANSACTION_STATUS_ENUM NOT NULL,
	type TRANSACTION_TYPE_ENUM NOT NULL,
	created_by TRANSACTION_CREATED_BY_ENUM NOT NULL,
	buyer_id UUID REFERENCES users NOT NULL,
	seller_id UUID REFERENCES users NOT NULL,
	delivery_duration INT NOT NULL,
	currency VARCHAR NOT NULL,
	charge_configuration JSON NOT NULL,
	product_details JSON NOT NULL,
	total_amount INT NOT NULL,
	total_cost INT NOT NULL,
	charges INT NOT NULL,
	receivable_amount INT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	deleted_at TIMESTAMPTZ,
	version INT DEFAULT 1,
);

CREATE TABLE IF NOT EXISTS transaction_timelines (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name TRANSACTION_TIMELINE_NAME_ENUM NOT NULL,
	transaction_id UUID REFERENCES transactions NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	deleted_at TIMESTAMPTZ,
	version INT DEFAULT 1,
);