BEGIN;

ALTER TABLE users
    DROP COLUMN role;

DROP TYPE user_role;

COMMIT;