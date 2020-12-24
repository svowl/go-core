package pg

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"filmoteka/pkg/films"
)

// Pg это структура БД (Postgres)
type Pg struct {
	connString string
}

// New возвращает новый объект Pg
func New(connString string) *Pg {
	var pg Pg
	pg.connString = connString
	return &pg
}

// connect возвращает контекст и соединение к БД
func (pg *Pg) connect() (context.Context, *pgxpool.Pool, error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(context.Background(), pg.connString)
	if err != nil {
		return nil, nil, err
	}
	return ctx, db, nil
}

// loadSQL выполняет загрузку SQL-запросов в виде текста в БД.
// Используется для загрузки тестового дампа в БД перед тестированием.
func (pg *Pg) loadSQL(sqlStr string) error {

	sqlStrings := strings.Split(sqlStr, ";\n")

	ctx, db, err := pg.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	batch := new(pgx.Batch)
	for _, sql := range sqlStrings {
		batch.Queue(sql)
	}
	res := tx.SendBatch(ctx, batch)
	err = res.Close()
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AddFilms добавляет фильмы в БД
func (pg *Pg) AddFilms(films []films.Film) error {
	ctx, db, err := pg.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	batch := new(pgx.Batch)
	for _, f := range films {
		batch.Queue(
			"INSERT INTO films(title, year, profit, pgrating, studio_id) VALUES ($1, $2, $3, $4, $5)",
			f.Title,
			f.Year,
			f.Profit,
			f.Pgrating,
			f.StudioID,
		)
	}
	res := tx.SendBatch(ctx, batch)
	err = res.Close()
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// DeleteFilm удаляет фильм из БД
func (pg *Pg) DeleteFilm(ID int) error {
	ctx, db, err := pg.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(ctx, `DELETE FROM films WHERE id=$1`, ID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateFilm обновляет данные фильма в БД
func (pg *Pg) UpdateFilm(ID int, film films.Film) error {
	ctx, db, err := pg.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(ctx, `
		UPDATE films
		SET title=$1, year=$2, profit=$3, pgrating=$4, studio_id=$5
		WHERE id=$6`,
		film.Title,
		film.Year,
		film.Profit,
		film.Pgrating,
		film.StudioID,
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Films возвращает один конкретный фильм или все фильмы из БД
func (pg *Pg) Films(studioID int) ([]films.Film, error) {
	ctx, db, err := pg.connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var rows pgx.Rows
	if studioID > 0 {
		rows, err = db.Query(ctx, `
			SELECT id, title, year, profit, pgrating, studio_id
			FROM films
			WHERE studio_id = $1`,
			studioID,
		)
	} else {
		rows, err = db.Query(ctx, `
			SELECT id, title, year, profit, pgrating, studio_id
			FROM films`,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []films.Film

	for rows.Next() {
		var f films.Film
		err := rows.Scan(
			&f.ID,
			&f.Title,
			&f.Year,
			&f.Profit,
			&f.Pgrating,
			&f.StudioID,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}
