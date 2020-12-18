-- actors - актеры
CREATE TABLE actors (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    year_of_birth INTEGER NOT NULL DEFAULT 0
);

-- directors - режиссеры
CREATE TABLE directors (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    year_of_birth INTEGER NOT NULL DEFAULT 0
);

-- studios - режиссеры
CREATE TABLE studios (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL DEFAULT ''
);

-- тип для поля рейтинга
CREATE TYPE pg_rating AS ENUM ('PG-10', 'PG-13', 'PG-18');

-- films - фильмы
CREATE TABLE films (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL, -- название
    year INTEGER DEFAULT 0 CHECK (year >= 1800), -- год выхода
    profit BIGINT DEFAULT 0, -- сборы
    pgrating pg_rating DEFAULT 'PG-10', -- рейтинг
    studio_id INTEGER REFERENCES studios(id) DEFAULT 0
);
-- уникальный индекс по названию фильма и году выхода
CREATE UNIQUE INDEX IF NOT EXISTS title_year_idx ON films USING btree (title, year);
-- на всякий случай: индекс для title не создаем, потому что
-- There's no need to manually create indexes on unique columns; doing so would just duplicate the automatically-created index.
-- https://www.postgresql.org/docs/current/indexes-unique.html

-- связь между фильмами и актерами
-- (у одного фильма может быть несколько актеров)
CREATE TABLE film_actors (
    id BIGSERIAL PRIMARY KEY,
    film_id BIGINT NOT NULL REFERENCES films(id),
    actor_id INTEGER NOT NULL REFERENCES actors(id),
    UNIQUE(film_id, actor_id)
);

-- связь между фильмами и режиссерами
-- (у одного фильма может быть несколько режиссеров)
CREATE TABLE film_directors (
    id BIGSERIAL PRIMARY KEY,
    film_id BIGINT NOT NULL REFERENCES films(id),
    director_id INTEGER NOT NULL REFERENCES directors(id),
    UNIQUE(film_id, director_id)
);