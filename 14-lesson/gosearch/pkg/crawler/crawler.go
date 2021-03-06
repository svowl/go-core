package crawler

// Поисковый робот.
// Осуществляет сканирование сайтов.

// Scanner определяет контракт поискового робота.
type Scanner interface {
	Scan(url []string, depth int) ([]Document, error)
}

// Document - документ, веб-страница, полученная поисковым роботом.
type Document struct {
	ID    int
	URL   string
	Title string
	Body  string
}
