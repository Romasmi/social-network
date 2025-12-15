CREATE TABLE IF NOT EXISTS cities (
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country_code CHAR(2) NOT NULL,
    state_province VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS IDX_cities_name_state_province_country_code ON cities(name, state_province, country_code);

CREATE INDEX idx_cities_country ON cities(country_code);
CREATE INDEX idx_cities_name ON cities(name);

-- For testing reason insert dummy data with Russian cities
INSERT INTO cities (id, name, country_code, state_province)
VALUES
    (gen_random_uuid(), 'Москва', 'RU', 'Московская область'),
    (gen_random_uuid(), 'Санкт-Петербург', 'RU', 'Ленинградская область'),
    (gen_random_uuid(), 'Новосибирск', 'RU', 'Новосибирская область'),
    (gen_random_uuid(), 'Екатеринбург', 'RU', 'Свердловская область'),
    (gen_random_uuid(), 'Нижний Новгород', 'RU', 'Нижегородская область'),
    (gen_random_uuid(), 'Казань', 'RU', 'Республика Татарстан'),
    (gen_random_uuid(), 'Челябинск', 'RU', 'Челябинская область'),
    (gen_random_uuid(), 'Омск', 'RU', 'Омская область'),
    (gen_random_uuid(), 'Самара', 'RU', 'Самарская область'),
    (gen_random_uuid(), 'Ростов-на-Дону', 'RU', 'Ростовская область'),
    (gen_random_uuid(), 'Уфа', 'RU', 'Республика Башкортостан'),
    (gen_random_uuid(), 'Красноярск', 'RU', 'Красноярский край'),
    (gen_random_uuid(), 'Пермь', 'RU', 'Пермский край'),
    (gen_random_uuid(), 'Воронеж', 'RU', 'Воронежская область'),
    (gen_random_uuid(), 'Волгоград', 'RU', 'Волгоградская область'),
    (gen_random_uuid(), 'Краснодар', 'RU', 'Краснодарский край'),
    (gen_random_uuid(), 'Саратов', 'RU', 'Саратовская область'),
    (gen_random_uuid(), 'Тюмень', 'RU', 'Тюменская область'),
    (gen_random_uuid(), 'Тольятти', 'RU', 'Самарская область'),
    (gen_random_uuid(), 'Ижевск', 'RU', 'Удмуртская Республика'),
    (gen_random_uuid(), 'Барнаул', 'RU', 'Алтайский край'),
    (gen_random_uuid(), 'Ульяновск', 'RU', 'Ульяновская область'),
    (gen_random_uuid(), 'Иркутск', 'RU', 'Иркутская область'),
    (gen_random_uuid(), 'Хабаровск', 'RU', 'Хабаровский край'),
    (gen_random_uuid(), 'Ярославль', 'RU', 'Ярославская область'),
    (gen_random_uuid(), 'Владивосток', 'RU', 'Приморский край'),
    (gen_random_uuid(), 'Махачкала', 'RU', 'Республика Дагестан'),
    (gen_random_uuid(), 'Томск', 'RU', 'Томская область'),
    (gen_random_uuid(), 'Оренбург', 'RU', 'Оренбургская область'),
    (gen_random_uuid(), 'Кемерово', 'RU', 'Кемеровская область');
