BEGIN;

CREATE type leave AS ENUM (
    'sick',
    'casual',
    'other'
);

CREATE TABLE user_leave (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    leave_from TIMESTAMP NOT NULL,
    leave_to TIMESTAMP NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    archived_at TIMESTAMP
);

COMMIT;