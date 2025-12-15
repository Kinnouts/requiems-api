CREATE TABLE IF NOT EXISTS words (
  id          SERIAL PRIMARY KEY,
  word        TEXT NOT NULL,
  definition  TEXT NOT NULL,
  part_of_speech TEXT
);

INSERT INTO words (word, definition, part_of_speech)
SELECT w, d, p FROM (VALUES
  ('requiem', 'A mass for the repose of the souls of the dead; a musical composition in honor of the dead.', 'noun'),
  ('latency', 'The delay before a transfer of data begins following an instruction.', 'noun'),
  ('idempotent', 'Producing the same result if applied multiple times.', 'adjective')
) AS seed(w, d, p)
WHERE NOT EXISTS (SELECT 1 FROM words);


