CREATE TABLE IF NOT EXISTS posts (
    id uuid NOT NULL PRIMARY KEY,
    profile_id uuid NOT NULL,
    text TEXT NOT NULL ,
    CONSTRAINT FK_profiles_user_id
        FOREIGN KEY (profile_id) REFERENCES profiles(id)
);
