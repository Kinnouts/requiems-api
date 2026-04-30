CREATE TABLE IF NOT EXISTS exercises (
    id                SERIAL PRIMARY KEY,
    external_id       VARCHAR(20)  NOT NULL UNIQUE,
    name              TEXT         NOT NULL,
    body_parts        TEXT[]       NOT NULL DEFAULT '{}',
    equipment         TEXT[]       NOT NULL DEFAULT '{}',
    target_muscles    TEXT[]       NOT NULL DEFAULT '{}',
    secondary_muscles TEXT[]       NOT NULL DEFAULT '{}',
    instructions      TEXT[]       NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_exercises_body_parts  ON exercises USING GIN(body_parts);
CREATE INDEX IF NOT EXISTS idx_exercises_equipment   ON exercises USING GIN(equipment);
CREATE INDEX IF NOT EXISTS idx_exercises_target      ON exercises USING GIN(target_muscles);
CREATE INDEX IF NOT EXISTS idx_exercises_secondary   ON exercises USING GIN(secondary_muscles);
CREATE INDEX IF NOT EXISTS idx_exercises_name_fts    ON exercises USING GIN(to_tsvector('english', name));
