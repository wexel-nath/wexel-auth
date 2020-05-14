
CREATE TABLE IF NOT EXISTS service (
    service_id          SERIAL      PRIMARY KEY NOT NULL,
    service_name        VARCHAR(16) NOT NULL,
    service_description TEXT        NOT NULL,
    UNIQUE (service_name)
);

INSERT INTO service (service_name, service_description)
VALUES
    ('all', 'Permissions for all services'),
    ('authentication', 'Handles authentication and permissions for all services'),
    ('continuance', 'CRM service for Continuance Pictures')
ON CONFLICT DO NOTHING;
