
CREATE TABLE IF NOT EXISTS user_permission (
    user_id       INTEGER NOT NULL REFERENCES users(user_id),
    permission_id INTEGER NOT NULL REFERENCES permission(permission_id),
    PRIMARY KEY (user_id, permission_id)
);
