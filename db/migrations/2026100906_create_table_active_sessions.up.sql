CREATE TABLE active_sessions (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    games_owned_id UUID NOT NULL,
    external_game_id VARCHAR(100) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT fk_session_must_be_owned 
        FOREIGN KEY (user_id, external_game_id) 
        REFERENCES games_owned(user_id, external_game_id) ON DELETE CASCADE,
    CONSTRAINT fk_session_owned_id_link 
        FOREIGN KEY (games_owned_id) 
        REFERENCES games_owned(id) ON DELETE CASCADE
);

CREATE INDEX idx_active_sessions_timeout ON active_sessions(user_id, expires_at);