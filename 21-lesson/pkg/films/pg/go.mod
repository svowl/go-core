module filmoteka/pkg/films/pg

go 1.15

replace filmoteka/pkg/films => ../

require (
	filmoteka/pkg/films v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v4 v4.10.1
)
