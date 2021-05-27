CREATE UNIQUE INDEX unique_user_email ON users(email) WHERE users.archived_at IS NULL;
CREATE UNIQUE INDEX unique_user_phone ON users(email) WHERE users.archived_at IS NULL AND users.phone IS NOT NULL;