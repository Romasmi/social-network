CREATE TABLE profiles
(
    id uuid PRIMARY KEY,
    user_id uuid UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    second_name VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,
    gender VARCHAR(10) DEFAULT NULL,
    biography TEXT,
    city_id uuid  NOT NULL,
    CONSTRAINT FK_profiles_user_id
        FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_profiles_city
        FOREIGN KEY (city_id) REFERENCES cities(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IDX_users_city ON profiles(city_id);
CREATE INDEX IDX_users_birthdate ON profiles(birthdate);
CREATE INDEX IDX_profile_name_composite ON profiles (first_name varchar_pattern_ops, second_name varchar_pattern_ops);
CREATE INDEX IDX_profile_last_name ON profiles (second_name varchar_pattern_ops);