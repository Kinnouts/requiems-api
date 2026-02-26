CREATE TABLE IF NOT EXISTS quotes (
  id      SERIAL PRIMARY KEY,
  text    TEXT NOT NULL,
  author  TEXT
);

INSERT INTO quotes (text, author)
SELECT v, a FROM (VALUES
  ('Simplicity is the soul of efficiency.', 'Austin Freeman'),
  ('Programs must be written for people to read, and only incidentally for machines to execute.', 'Harold Abelson'),
  ('Premature optimization is the root of all evil.', 'Donald Knuth')
) AS seed(v, a)
WHERE NOT EXISTS (SELECT 1 FROM quotes);


