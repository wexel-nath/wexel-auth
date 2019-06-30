
CREATE TABLE IF NOT EXISTS role (
    role_id          SERIAL      PRIMARY KEY NOT NULL,
    role_name        VARCHAR(16) NOT NULL,
    role_description TEXT        NOT NULL,
    service_id       INTEGER     NOT NULL REFERENCES service(service_id)
);
