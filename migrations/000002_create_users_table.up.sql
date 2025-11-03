CREATE TABLE users
(
    id uuid PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);


-- Note: indexes are not created intentionally to get performance issue
-- CREATE INDEX IDX_users_city ON users(city);
-- CREATE INDEX IDX_users_birthdate ON users(birthdate);
-- CREATE INDEX IDX_users_name ON users(name);
