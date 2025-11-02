CREATE TABLE cities (
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country_code CHAR(2) NOT NULL,
    state_province VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE UNIQUE INDEX IDX_cities_name_state_province_country_code ON cities(name, state_province, country_code);

-- Note: indexes are not created intentionally to get performance issues
-- CREATE INDEX idx_cities_country ON cities(country_code);
-- CREATE INDEX idx_cities_name ON cities(name);