INSERT INTO actors (id, first_name) VALUES (0, 'actor not specified');
ALTER SEQUENCE actors_id_seq RESTART WITH 100;

INSERT INTO directors (id, first_name) VALUES (0, 'director not specified');
ALTER SEQUENCE directors_id_seq RESTART WITH 100;

INSERT INTO studios (id, name) VALUES (0, 'studio not specified');
ALTER SEQUENCE studios_id_seq RESTART WITH 100;
