
CREATE TABLE IF NOT EXISTS user (
    user_id    SERIAL      PRIMARY KEY NOT NULL,
    first_name VARCHAR(16) NOT NULL,
    last_name  VARCHAR(16) NOT NULL,
    email      VARCHAR(32) NOT NULL,
    username   VARCHAR(32) NOT NULL,
    password   VARCHAR(32) NOT NULL
);
