-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    avatar VARCHAR(255),
    bio TEXT,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create a default admin user with a proper Argon2id hash for password 'admin'
INSERT INTO users (username, password, email, first_name, is_admin) 
VALUES ('admin', '$argon2id$v=19$m=65536,t=1,p=4$DqHRSyN6rT1FshvOM5wchg$5RLPolMaBW10ooTkE3ZtjOJ3VlRiyXcIcBgV6HPcFXQ', 'admin@example.com', 'Admin', TRUE);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users; 