CREATE TABLE IF NOT EXISTS commodity_price_history (
    slug       VARCHAR(50)   NOT NULL,
    name       TEXT          NOT NULL,
    unit       VARCHAR(50)   NOT NULL,
    currency   CHAR(3)       NOT NULL DEFAULT 'USD',
    year       SMALLINT      NOT NULL,
    price      NUMERIC(14,4) NOT NULL,
    source     TEXT          NOT NULL DEFAULT 'fred',
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    PRIMARY KEY (slug, year)
);

CREATE INDEX IF NOT EXISTS idx_commodity_price_history_slug_year
    ON commodity_price_history (slug, year DESC);
