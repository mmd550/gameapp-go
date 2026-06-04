ALTER TABLE users ADD COLUMN password TEXT;

DELETE FROM users WHERE password is null;

ALTER TABLE users ALTER COLUMN password SET NOT NULL;