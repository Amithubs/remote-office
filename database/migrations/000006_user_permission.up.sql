CREATE type permission AS ENUM (
    'employee',
    'admin'
    );
CREATE TABLE user_permission(
    user_id int NOT NULL REFERENCES users(id),
    permission_type permission NOT NULL DEFAULT 'employee'::permission,
    UNIQUE(user_id,permission_type)
);