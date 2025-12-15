CREATE TABLE IF NOT EXISTS advice (
  id   SERIAL PRIMARY KEY,
  text TEXT NOT NULL
);

INSERT INTO advice (text)
SELECT v FROM (VALUES
  ('Ship small, ship often.'),
  ('Talk to users before you over-engineer.'),
  ('Good logs today save you hours tomorrow.'),
  ('Automate anything you do more than twice.'),
  ('Optimize for readability over cleverness.'),
  ('Start with the simplest data model that could work.')
) AS seed(v)
WHERE NOT EXISTS (SELECT 1 FROM advice);


