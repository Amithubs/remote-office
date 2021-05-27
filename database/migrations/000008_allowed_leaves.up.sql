CREATE TABLE allowed_leaves(
    id SERIAL PRIMARY KEY ,
    year int NOT NULL,
    leave_type leave NOT NULL ,
    allowed_leaves int NOT NULL,
    UNIQUE (leave_type,allowed_leaves)
);