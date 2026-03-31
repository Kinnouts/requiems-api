CREATE TABLE IF NOT EXISTS inflation_data (
    country_code  CHAR(2)       NOT NULL,
    country_name  TEXT          NOT NULL DEFAULT '',
    year          SMALLINT      NOT NULL,
    rate          NUMERIC(8,4)  NOT NULL,
    source        TEXT          NOT NULL DEFAULT 'world_bank',
    last_updated  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    PRIMARY KEY   (country_code, year)
);

CREATE INDEX IF NOT EXISTS idx_inflation_data_country_year
    ON inflation_data (country_code, year DESC);
