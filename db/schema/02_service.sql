
CREATE TABLE IF NOT EXISTS service (
    service_id          SERIAL      PRIMARY KEY NOT NULL,
    service_name        VARCHAR(16) NOT NULL,
    service_description TEXT        NOT NULL
);
