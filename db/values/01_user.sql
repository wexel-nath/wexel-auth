
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
