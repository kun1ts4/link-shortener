-- +goose Up
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    original TEXT NOT NULL UNIQUE,
    short VARCHAR(20) NOT NULL UNIQUE,
    clicks BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS links_short ON links(short);
CREATE INDEX IF NOT EXISTS links_original ON links(original);

-- +goose Down
SELECT 'down SQL query';

DROP TABLE IF EXISTS links;