CREATE TABLE IF NOT EXISTS profile_friends (
    profile1_id uuid NOT NULL,
    profile2_id uuid NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT FK_profiles_user_id
        FOREIGN KEY (profile1_id) REFERENCES profiles(id),
    CONSTRAINT FK_profiles_friend_id
        FOREIGN KEY (profile2_id) REFERENCES profiles(id),
    PRIMARY KEY (profile1_id, profile2_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS UX_profile_friends_profile_id ON profile_friends (
                                                           LEAST(profile1_id, profile2_id),
                                                           GREATEST(profile1_id, profile2_id)
);

CREATE INDEX IF NOT EXISTS IDX_user_friends_profile2_id ON profile_friends (profile2_id);
