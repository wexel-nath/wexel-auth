
CREATE TABLE IF NOT EXISTS session (
    session_id      TEXT      PRIMARY KEY NOT NULL,
    user_id         INTEGER   REFERENCES users (user_id),
    session_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    session_expiry  TIMESTAMP WITH TIME ZONE NOT NULL
);
