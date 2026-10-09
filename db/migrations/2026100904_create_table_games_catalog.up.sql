CREATE TABLE games_catalog (
    external_game_id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    price_in_balance BIGINT NOT NULL,
    is_available BOOLEAN DEFAULT TRUE NOT NULL
);