
INSERT INTO service (service_name, service_description)
VALUES
    ('all', 'Permissions for all services'),
    ('authentication', 'Handles authentication and permissions for all services'),
    ('continuance', 'CRM service for Continuance Pictures')
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

SELECT add_permission('user.create', 'Create a user', 'all');
SELECT add_permission('user.edit', 'Edit another user', 'all');
SELECT add_permission('user.delete', 'Delete another user', 'all');
SELECT add_permission('continuance', 'Base level permissions for continuance', 'continuance');

DROP FUNCTION IF EXISTS add_permission(name TEXT, description TEXT, s_name TEXT);

CREATE OR REPLACE FUNCTION add_user_permission(u_name TEXT, p_name TEXT) RETURNS VOID AS '
    BEGIN
        INSERT INTO user_permission (user_id, permission_id)
        SELECT user_id, permission_id
        FROM   users, permission
        WHERE  username = u_name
        AND    permission_name = p_name;
    END;
'
LANGUAGE 'plpgsql';

SELECT add_user_permission('admin', 'user.create');
SELECT add_user_permission('admin', 'user.edit');
SELECT add_user_permission('admin', 'user.delete');
SELECT add_user_permission('admin', 'continuance');

DROP FUNCTION IF EXISTS add_user_permission(username TEXT, permission_name TEXT);
