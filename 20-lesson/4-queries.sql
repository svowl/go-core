-- выборка фильмов с названием студии; 
SELECT f.*, s.name
FROM films f
JOIN studios s ON f.studio_id = s.id 
WHERE s.name = 'Paramount Pictures';

-- выборка фильмов для некоторого актёра; 
SELECT f.*, a.first_name, a.last_name 
FROM films f 
JOIN film_actors fa ON fa.film_id = f.id 
JOIN actors a ON fa.actor_id = a.id
WHERE last_name='Cruise';

-- подсчёт фильмов для некоторого режиссёра;
SELECT d.first_name, d.last_name, COUNT(f.id)
FROM films f
JOIN film_directors fd ON fd.film_id = f.id
JOIN directors d ON fd.director_id = d.id
WHERE d.last_name = 'Tarantino'
GROUP BY d.first_name, d.last_name;

-- выборка фильмов для нескольких режиссёров из списка (подзапрос); 
SELECT f.*
FROM films f
WHERE f.id IN (
    SELECT film_id
    FROM film_directors
    WHERE director_id IN (
        SELECT id
        FROM directors
        WHERE last_name IN ('Tarantino', 'Kubrick')
    )
);

-- то же самое, другой вариант
WITH director_ids AS (
    SELECT id
    FROM directors
    WHERE last_name IN ('Tarantino', 'Kubrick')
), film_ids AS (
    SELECT film_id
    FROM film_directors
    WHERE director_id IN (SELECT id FROM director_ids)
)
SELECT f.*
FROM films f
WHERE f.id IN (SELECT film_id FROM film_ids);

-- то же самое без подзапроса
SELECT f.*, d.first_name, d.last_name
FROM films f
JOIN film_directors fd ON fd.film_id = f.id
JOIN directors d ON fd.director_id = d.id
WHERE d.last_name = 'Tarantino' OR d.last_name = 'Kubrick';

-- подсчёт количества фильмов для актёра; 
SELECT a.first_name, a.last_name, COUNT(f.id)
FROM films f
JOIN film_actors fa ON fa.film_id = f.id
JOIN actors a ON fa.actor_id = a.id
WHERE a.last_name = 'Cruise'
GROUP BY a.first_name, a.last_name;

-- выборка актёров и режиссёров, участвовавших более чем в 2 фильмах; 
SELECT first_name, last_name, cnt
FROM (
    SELECT a.first_name, a.last_name, COUNT(f.id) AS cnt
    FROM films f
    JOIN film_actors fa ON fa.film_id = f.id
    JOIN actors a ON fa.actor_id = a.id
    GROUP BY a.first_name, a.last_name
    HAVING COUNT(f.id) > 2
) a
UNION
(
    SELECT d.first_name, d.last_name, COUNT(f.id) AS cnt
    FROM films f
    JOIN film_directors fd ON fd.film_id = f.id
    JOIN directors d ON fd.director_id = d.id
    GROUP BY d.first_name, d.last_name
    HAVING COUNT(f.id) > 2
);

-- подсчёт количества фильмов со сборами больше $400млн;
SELECT COUNT(title)
FROM (
    SELECT title FROM films WHERE profit > 400000000
) f;

-- то же самое, но с выводом этих фильмов
SELECT title, profit FROM films WHERE profit > 400000000;

-- подсчитать количество режиссёров, фильмы которых собрали больше $400млн; 
SELECT d.first_name, d.last_name, SUM(f.profit)
FROM directors d
JOIN film_directors fd ON fd.director_id = d.id
JOIN films f ON fd.film_id = f.id
GROUP BY d.first_name, d.last_name
HAVING SUM(f.profit) > 400000000
ORDER BY SUM(f.profit) DESC;

-- выборка различных фамилий актёров; 
SELECT DISTINCT(last_name) FROM actors;
SELECT last_name FROM actors GROUP BY last_name;

-- подсчёт количества фильмов, имеющих дубли по названию.

-- сначала создадим дубли:
INSERT INTO films (id, title, year, profit, pgrating, studio_id)
VALUES
    (100, 'Pretty Woman', 1991, 463406268, 'PG-13', 1),
    (300, 'Notting Hill', 1998, 364000000, 'PG-13', 3);

SELECT COUNT(title)
FROM (
	SELECT f.title
	FROM films f
	GROUP BY f.title
	HAVING COUNT(f.title) > 1
) f;

-- то же самое, но с выводом дублей
SELECT f.title, COUNT(f.title)
FROM films f
GROUP BY f.title
HAVING COUNT(f.title) > 1;

