package films

// Film это структура, описывающая фильм
type Film struct {
	ID       int
	Title    string
	Year     int
	Profit   int
	Pgrating string
	StudioID int
}

// Studio это структура, описывающая студию
type Studio struct {
	ID   int
	Name string
}

// IFilms это интерфейс для работы с фильмами
type IFilms interface {
	AddFilms([]Film) error
	DeleteFilm(int) error
	UpdateFilm(int, Film) error
	Films(int) ([]Film, error)
}
