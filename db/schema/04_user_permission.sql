
CREATE TABLE IF NOT EXISTS user_permission (
    user_id       INTEGER NOT NULL REFERENCES users(user_id),
    permission_id INTEGER NOT NULL REFERENCES permission(permission_id),
    PRIMARY KEY (user_id, permission_id)
);

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
