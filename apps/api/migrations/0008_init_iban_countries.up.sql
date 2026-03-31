CREATE TABLE iban_countries (
    country_code    CHAR(2)      NOT NULL PRIMARY KEY,
    country_name    TEXT         NOT NULL,
    iban_length     SMALLINT     NOT NULL,
    bban_format     TEXT         NOT NULL DEFAULT '',
    bank_offset     SMALLINT     NOT NULL DEFAULT 0,
    bank_length     SMALLINT     NOT NULL DEFAULT 0,
    account_offset  SMALLINT     NOT NULL DEFAULT 0,
    account_length  SMALLINT     NOT NULL DEFAULT 0,
    sepa_member     BOOLEAN      NOT NULL DEFAULT false,
    last_updated    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Fast lookup for SEPA-only queries.
CREATE INDEX iban_countries_sepa_idx ON iban_countries (sepa_member)
    WHERE sepa_member = true;
