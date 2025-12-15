CREATE TABLE user_friends (
    user1_id uuid UNIQUE,
    user2_id uuid NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT FK_profiles_user_id
        FOREIGN KEY (user1_id) REFERENCES users(id),
    CONSTRAINT FK_profiles_friend_id
        FOREIGN KEY (user2_id) REFERENCES users(id),
    PRIMARY KEY (user1_id, user2_id)
);

CREATE UNIQUE INDEX idx_unique_friendships ON user_friends (
                                                           LEAST(user1_id, user2_id),
                                                           GREATEST(user1_id, user2_id)
);

CREATE INDEX IDX_user_friends_user1_id ON user_friends(user1_id);
CREATE INDEX IDX_user_friends_user2_id ON user_friends(user2_id);
