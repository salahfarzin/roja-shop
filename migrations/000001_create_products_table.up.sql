-- +migrate Up
CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY,
    brand TEXT NOT NULL,
    title TEXT NOT NULL,
    inventory INTEGER NOT NULL,
    price DOUBLE NOT NULL,
    old_price DOUBLE,
    discount DOUBLE,
    description TEXT NOT NULL,
    details TEXT,       -- JSON string
    style_notes TEXT    -- JSON string
);
