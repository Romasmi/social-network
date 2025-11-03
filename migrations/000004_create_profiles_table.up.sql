CREATE TYPE gender AS ENUM ('male', 'female');
CREATE TABLE profiles
(
    id uuid PRIMARY KEY,
    user_id uuid UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    second_name VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,
    gender gender NOT NULL,
    biography TEXT,
    city uuid  NOT NULL,
    CONSTRAINT FK_profiles_user_id
        FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_profiles_city
        FOREIGN KEY (city) REFERENCES cities(id),
    created_at TIMESTAMP DEFAULT NOW()
);


-- Note: indexes are not created intentionally to get performance issue
-- CREATE INDEX IDX_users_city ON users(city);
-- CREATE INDEX IDX_users_birthdate ON users(birthdate);
-- CREATE INDEX IDX_users_name ON users(name);
