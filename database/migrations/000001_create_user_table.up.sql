BEGIN;

CREATE TYPE designation AS ENUM (
    'ceo',
    'hr',
    'ui',
    'sde-1',
    'sde-2',
    'team-lead',
    'trainee'
);

CREATE TYPE image_type AS ENUM (
    'profile'
);

CREATE TABLE images (
    id SERIAL primary key,
    type image_type NOT NULL,
    bucket TEXT,/*??*/
    path TEXT,
    created_at TIMESTAMP DEFAULT now(),
    archived_at TIMESTAMP
);

CREATE TABLE users (
   id SERIAL primary key,
   name TEXT NOT NULL,
   phone TEXT,
   email TEXT NOT NULL,
   password TEXT NOT NULL,
   position designation DEFAULT 'trainee'::designation,
   profile_image INT REFERENCES images(id),
   created_at TIMESTAMP DEFAULT now(),
   archived_at TIMESTAMP
);

COMMIT;