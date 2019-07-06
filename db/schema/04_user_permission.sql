
CREATE TABLE IF NOT EXISTS permission (
    permission_id          SERIAL      PRIMARY KEY NOT NULL,
    permission_name        VARCHAR(16) NOT NULL,
    permission_description TEXT        NOT NULL,
    service_id             INTEGER     NOT NULL REFERENCES service(service_id)
);
