ALTER TABLE tokens DROP COLUMN hash;
ALTER TABLE tokens ADD COLUMN hash BYTEA NOT NULL;