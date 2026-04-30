CREATE TABLE IF NOT EXISTS bin_data (
    bin_prefix     VARCHAR(8)    PRIMARY KEY,
    prefix_length  SMALLINT      NOT NULL,
    scheme         TEXT          NOT NULL DEFAULT '',
    card_type      TEXT          NOT NULL DEFAULT '',  -- credit | debit | prepaid | charge
    card_level     TEXT          NOT NULL DEFAULT '',  -- classic | gold | platinum | business | corporate | signature | infinite | standard
    issuer_name    TEXT          NOT NULL DEFAULT '',
    issuer_url     TEXT          NOT NULL DEFAULT '',
    issuer_phone   TEXT          NOT NULL DEFAULT '',
    country_code   CHAR(2)       NOT NULL DEFAULT '',  -- ISO 3166-1 alpha-2
    country_name   TEXT          NOT NULL DEFAULT '',
    prepaid        BOOLEAN       NOT NULL DEFAULT FALSE,
    source         TEXT          NOT NULL DEFAULT '',  -- which dataset(s) provided this row
    confidence     NUMERIC(3,2)  NOT NULL DEFAULT 0.50, -- 0.00–1.00
    first_seen     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    last_updated   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Fast fallback for 6-digit prefix when an 8-digit BIN is not found
CREATE INDEX IF NOT EXISTS idx_bin_data_prefix6
    ON bin_data (LEFT(bin_prefix, 6));

-- Analytics / filtering indexes
CREATE INDEX IF NOT EXISTS idx_bin_data_scheme
    ON bin_data (scheme);

CREATE INDEX IF NOT EXISTS idx_bin_data_country
    ON bin_data (country_code);
