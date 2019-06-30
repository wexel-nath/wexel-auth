
CREATE TABLE IF NOT EXISTS permission (
    permission_id SERIAL  PRIMARY KEY NOT NULL,
    user_id       INTEGER NOT NULL REFERENCES user(user_id),
    role_id       INTEGER NOT NULL REFERENCES role(role_id)
);
