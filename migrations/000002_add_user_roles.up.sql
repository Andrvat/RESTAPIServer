BEGIN;

CREATE TYPE user_role AS ENUM (
    'basic',
    'admin',
    'moderator'
    );

ALTER TABLE users
    ADD COLUMN role user_role DEFAULT 'basic';

UPDATE users
SET role='basic';

COMMIT;
