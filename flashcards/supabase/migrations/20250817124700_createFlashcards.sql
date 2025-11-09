CREATE SCHEMA IF NOT EXISTS flashcards;

CREATE TABLE IF NOT EXISTS flashcards.flashcards (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_flashcards_created_at ON flashcards.flashcards(created_at);