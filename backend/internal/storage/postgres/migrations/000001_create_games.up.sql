CREATE TABLE IF NOT EXISTS games (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    genre TEXT NOT NULL,
    platform TEXT NOT NULL,
    release_year INTEGER NOT NULL,
    rating INTEGER NOT NULL,
    status TEXT NOT NULL
);
