DROP TABLE IF EXISTS film_actors;
DROP TABLE IF EXISTS film_directors;
DROP TABLE IF EXISTS films;
DROP TABLE IF EXISTS studios;
DROP TABLE IF EXISTS directors;
DROP TABLE IF EXISTS actors;

DROP TYPE IF EXISTS pg_rating;

-- actors - актеры
CREATE TABLE actors (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    year_of_birth INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS actors_last_name_idx ON actors USING btree (lower(last_name));

-- directors - режиссеры
CREATE TABLE directors (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    year_of_birth INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS directors_last_name_idx ON directors USING btree (lower(last_name));

-- studios - режиссеры
CREATE TABLE studios (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS name_idx ON studios USING btree (lower(name));

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
CREATE UNIQUE INDEX IF NOT EXISTS title_year_idx ON films USING btree (lower(title), year);
-- на всякий случай: индекс для title не создаем, потому что
-- There's no need to manually create indexes on unique columns; doing so would just duplicate the automatically-created index.
-- https://www.postgresql.org/docs/current/indexes-unique.html
CREATE INDEX IF NOT EXISTS profit_idx ON films USING btree (profit);
CREATE INDEX IF NOT EXISTS year_idx ON films USING btree (year);

-- связь между фильмами и актерами
-- (у одного фильма может быть несколько актеров)
CREATE TABLE film_actors (
    film_id BIGINT NOT NULL REFERENCES films(id),
    actor_id INTEGER NOT NULL REFERENCES actors(id),
    UNIQUE(film_id, actor_id)
);

-- связь между фильмами и режиссерами
-- (у одного фильма может быть несколько режиссеров)
CREATE TABLE film_directors (
    film_id BIGINT NOT NULL REFERENCES films(id),
    director_id INTEGER NOT NULL REFERENCES directors(id),
    UNIQUE(film_id, director_id)
);