CREATE TABLE sessions
(
    id uuid PRIMARY KEY,
    user_id uuid,
    metadata JSONB NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT FK_profiles_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
);

