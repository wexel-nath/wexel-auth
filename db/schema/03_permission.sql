
CREATE TABLE IF NOT EXISTS permission (
    permission_id          SERIAL      PRIMARY KEY NOT NULL,
    permission_name        VARCHAR(16) NOT NULL,
    permission_description TEXT        NOT NULL,
    service_id             INTEGER     NOT NULL REFERENCES service(service_id),
    UNIQUE (service_id, permission_name)
);

CREATE OR REPLACE FUNCTION add_permission(name TEXT, description TEXT, s_name TEXT) RETURNS VOID AS '
    BEGIN
        INSERT INTO permission (permission_name, permission_description, service_id)
        SELECT name, description, service_id
        FROM service
        WHERE service_name = s_name;
    END;
'
LANGUAGE 'plpgsql';

SELECT add_permission('user.create', 'Create a user', 'all');
SELECT add_permission('user.edit', 'Edit another user', 'all');
SELECT add_permission('user.delete', 'Delete another user', 'all');
SELECT add_permission('continuance', 'Base level permissions for continuance', 'continuance');

DROP FUNCTION IF EXISTS add_permission(name TEXT, description TEXT, s_name TEXT);
