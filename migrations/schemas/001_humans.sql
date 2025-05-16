-- +goose Up
CREATE TABLE IF NOT EXISTS humans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT,
    age INT NOT NULL,
    gender TEXT NOT NULL,
    country TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS humans;