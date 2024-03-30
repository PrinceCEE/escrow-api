CREATE TYPE MODEL_STATUS_ENUM AS ENUM (
  "Successful",
  "Canceled",
  "Pending"
);

CREATE TYPE WITHDRAWAL_TYPE_ENUM AS ENUM (
  "Withdrawal",
  "Deposit"
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

