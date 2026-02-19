CREATE TABLE IF NOT EXISTS counters (
    namespace  TEXT      PRIMARY KEY,
    total      BIGINT    NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL
);
