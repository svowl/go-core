INSERT INTO actors (id, first_name, last_name, year_of_birth)
VALUES
    (1, 'Johnny', 'Depp', 1963),
    (2, 'Leonardo', 'DiCaprio', 1974),
    (3, 'Brad', 'Pitt', 1963),
    (4, 'Tom', 'Cruise', 1962),
    (5, 'Nicole', 'Kidman', 1967),
    (6, 'Hugh', 'Grant', 1960),
    (7, 'Julia', 'Roberts', 1967),
    (8, 'Richard', 'Gere', 1949),
    (9, 'Jean', 'Reno', 1948);
ALTER SEQUENCE actors_id_seq RESTART WITH 100;

INSERT INTO directors (id, first_name, last_name, year_of_birth)
VALUES
    (1, 'Garry', 'Marshall', 1934),
    (2, 'P.J.', 'Hogan', 1962),
    (3, 'Roger','Michell', 1956),
    (4, 'Stanley', 'Kubrick', 1928),
    (5, 'Brian', 'De Palma', 1940),
    (6, 'Neil', 'Jordan', 1950),
    (7, 'Quentin', 'Tarantino', 1963),
    (8, 'Lasse', 'Hallström', 1946),
    (9, 'Gore', 'Verbinski', 1964);
ALTER SEQUENCE directors_id_seq RESTART WITH 100;

INSERT INTO studios (id, name)
VALUES
    (1, 'Touchstone Pictures'),
    (2, 'Sony Pictures'),
    (3, 'Universal Pictures'),
    (4, 'Warner Bros.'),
    (5, 'Paramount Pictures'),
    (6, 'Geffen Pictures'),
    (7, 'Columbia Pictures'),
    (8, 'Walt Disney Pictures');
ALTER SEQUENCE studios_id_seq RESTART WITH 100;

INSERT INTO films (id, title, year, profit, pgrating, studio_id)
VALUES
    (1, 'Pretty Woman', 1990, 463406268, 'PG-13', 1),
    (2, 'My Best Friend''s Wedding', 1997, 290300000, 'PG-10', 2),
    (3, 'Notting Hill', 1999, 364000000, 'PG-13', 3),
    (4, 'Eyes Wide Shut', 1999, 162000000, 'PG-18', 4),
    (5, 'Mission: Impossible', 1996, 456494803, 'PG-18', 5),
    (6, 'Interview with the Vampire: The Vampire Chronicles', 1994, 223664608, 'PG-18', 6),
    (7, 'Once Upon a Time... in Hollywood', 2019, 374343626, 'PG-13', 7),
    (8, 'Inglourious Basterds', 2009, 321455689, 'PG-18', 3),
    (9, 'What''s Eating Gilbert Grape', 1993, 10000000, 'PG-10', 5),
    (10, 'Pirates of the Caribbean: The Curse of The Black Pearl', 2003, 654264015, 'PG-13', 8);
ALTER SEQUENCE films_id_seq RESTART WITH 100;

INSERT INTO film_actors (film_id, actor_id)
VALUES
    (1, 7),
    (1, 8),
    (2, 7),
    (3, 6),
    (3, 7),
    (4, 4),
    (4, 5),
    (5, 4),
    (5, 9),
    (6, 4),
    (6, 3),
    (7, 3),
    (7, 2),
    (8, 3),
    (9, 1),
    (9, 2),
    (10, 1);

INSERT INTO film_directors (film_id, director_id)
VALUES
    (1, 1),
    (2, 2),
    (3, 3),
    (4, 4),
    (5, 5),
    (6, 6),
    (7, 7),
    (7, 5),
    (8, 7),
    (9, 8),
    (10, 9);
    