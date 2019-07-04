
CREATE TABLE IF NOT EXISTS session (
    session_id TEXT        PRIMARY KEY NOT NULL,
    user_id    INTEGER     REFERENCES users(user_id),
    timestamp  INTEGER     NOT NULL,
    expiry     INTEGER     NOT NULL
);
