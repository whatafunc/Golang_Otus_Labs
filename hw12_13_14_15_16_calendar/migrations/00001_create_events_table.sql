-- +goose Up
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(40) NOT NULL,
    description TEXT NOT NULL,
    start TIMESTAMP,
    "end" TIMESTAMP,
    allday FLOAT NOT NULL,
    clinic TEXT,
    userid INT,
    service TEXT
);

-- +goose Down
DROP TABLE IF EXISTS events;