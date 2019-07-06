
INSERT INTO service (service_name, service_description)
VALUES
    ('ALL', 'Permissions for all services'),
    ('Authentication', 'Handles authentication and permissions for all services'),
    ('Continuance', 'CRM service for Continuance Pictures')
ON CONFLICT DO NOTHING;

CREATE OR REPLACE FUNCTION add_permission(name TEXT, description TEXT, s_name TEXT) RETURNS VOID AS '
    BEGIN
        INSERT INTO permission (permission_name, permission_description, service_id)
        SELECT name, description, service_id
        FROM service
        WHERE service_name = s_name;
    END;
'
LANGUAGE 'plpgsql';

SELECT add_permission('user.create', 'Create a user', 'ALL');
SELECT add_permission('user.view', 'View a user', 'ALL');

DROP FUNCTION IF EXISTS add_permission(name TEXT, description TEXT, s_name TEXT);
