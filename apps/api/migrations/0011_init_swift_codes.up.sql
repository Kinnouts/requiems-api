CREATE TABLE swift_codes (
    swift_code    VARCHAR(11)  NOT NULL PRIMARY KEY, -- full 11-char BIC, always uppercase
    bank_code     VARCHAR(4)   NOT NULL,             -- chars 1-4, institution code
    country_code  VARCHAR(2)   NOT NULL,             -- chars 5-6, ISO 3166-1 alpha-2
    location_code VARCHAR(2)   NOT NULL,             -- chars 7-8, location code
    branch_code   VARCHAR(3)   NOT NULL DEFAULT 'XXX', -- chars 9-11; 'XXX' = primary office
    bank_name     TEXT         NOT NULL DEFAULT '',
    city          TEXT         NOT NULL DEFAULT '',
    country_name  TEXT         NOT NULL DEFAULT '',
    last_updated  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Fast lookup for all branches of a bank in a country.
CREATE INDEX swift_codes_bank_country_idx ON swift_codes (bank_code, country_code);

-- Fast lookup for all codes by country.
CREATE INDEX swift_codes_country_idx ON swift_codes (country_code);
