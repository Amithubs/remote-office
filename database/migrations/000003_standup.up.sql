CREATE TABLE standup (
    id SERIAL primary key,
    user_id INT REFERENCES users(id),
    data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    archived_at TIMESTAMP
);