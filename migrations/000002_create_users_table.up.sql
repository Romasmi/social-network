CREATE TYPE gender AS ENUM ('male', 'female');
CREATE TABLE users
(
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,
    gender gender NOT NULL,
    interests TEXT,
    city uuid  NOT NULL,
    CONSTRAINT FK_users_city
    FOREIGN KEY (city) REFERENCES cities(id),
    created_at TIMESTAMP DEFAULT NOW()
);


-- Note: indexes are not created intentionally to get performance issue
-- CREATE INDEX IDX_users_city ON users(city);
-- CREATE INDEX IDX_users_birthdate ON users(birthdate);
-- CREATE INDEX IDX_users_name ON users(name);
