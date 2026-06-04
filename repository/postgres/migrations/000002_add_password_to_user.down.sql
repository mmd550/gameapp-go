-- Reverse: Remove NOT NULL constraint
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;

-- Reverse: Drop the column
ALTER TABLE users DROP COLUMN password;