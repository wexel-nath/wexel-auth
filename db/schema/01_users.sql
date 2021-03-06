
CREATE TABLE IF NOT EXISTS users (
    user_id    SERIAL      PRIMARY KEY NOT NULL,
    first_name VARCHAR(16) NOT NULL,
    last_name  VARCHAR(16) NOT NULL,
    email      VARCHAR(32) NOT NULL,
    username   VARCHAR(32) NOT NULL UNIQUE,
    password   TEXT        NOT NULL
);

CREATE EXTENSION pgcrypto;

INSERT INTO users(
  first_name,
  last_name,
  email,
  username,
  password
)
VALUES (
  'Admin',
  'User',
  'admin@getwexel.com',
  'admin',
  crypt('admin', gen_salt('bf'))
);
