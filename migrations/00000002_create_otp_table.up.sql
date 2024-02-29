CREATE TYPE OTP_TYPE AS ENUM ('SMS', 'EMAIL');

CREATE TABLE IF NOT EXISTS otps (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users,
  code CHAR(4) NOT NULL,
  is_used BOOLEAN NOT NULL,
  otp_type OTP_TYPE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ,
  version INT DEFAULT 1
);