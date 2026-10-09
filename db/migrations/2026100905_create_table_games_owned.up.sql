CREATE TABLE games_owned (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    external_game_id VARCHAR(100) NOT NULL REFERENCES games_catalog(external_game_id),
    acquired_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    playtime_minutes INT DEFAULT 0 NOT NULL,
    CONSTRAINT unique_user_game UNIQUE (user_id, external_game_id)
);

CREATE INDEX idx_games_owned_user ON games_owned(user_id);