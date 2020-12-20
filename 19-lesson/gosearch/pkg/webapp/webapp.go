package webapp

import (
	"gosearch/pkg/index"
	"gosearch/pkg/storage"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Service это служба Web-приложения, содержит ссылки на объекты роутера, БД и индекса
type Service struct {
	router *mux.Router
	db     *storage.Db
	index  *index.Index
}

// New создает объект Service, объявляет endpoints
func New(router *mux.Router, db *storage.Db, index *index.Index) *Service {
	var s Service
	s.router = router
	s.db = db
	s.index = index

	s.endpoints()

	return &s
}

// Определяем endpoints
func (s *Service) endpoints() {
	s.router.HandleFunc("/index", s.indexHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/docs", s.docsHandler).Methods(http.MethodGet)
}

// HTTP-обработчик /index выводит в HTML-формате содержимое индекса
func (s *Service) indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("main")
	tpl := `
<!DOCTYPE html>
<html>
<head>
<meta content="text/html; charset=utf-8" http-equiv="content-type">
<title>Gosearch: Index</title>
</head>
<body>
<h1>Gosearch index</h1>
<h3>Total rows: {{len .}}</h3>
<ul style="column-count: 4;">
{{range $word, $ids := .}}
<li>{{$word}}
{{end}}
</ul>
</body>
</html>
`
	t, err := t.Parse(tpl)
	if err != nil {
		http.Error(w, "ошибка при обработке шаблона", http.StatusInternalServerError)
		return
	}

	t.Execute(w, s.index.Hash)
}

// HTTP-обработчик /docs выводит в HTML-формате список документов из БД
func (s *Service) docsHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("main")
	tpl := `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=UTF-8">
<title>Gosearch: Documents</title>
</head>
<body>
<h1>Gosearch documents</h1>
<h3>Total documents: {{len .}}</h3>
<ul style="column-count: 3;">
{{range $doc := .}}
<li><a href="{{$doc.URL}}">{{if $doc.Title}}{{$doc.Title}}{{else}}&lt;No title&gt;{{end}}</a>
{{end}}
</ul>
</body>
</html>
`
	t, err := t.Parse(tpl)
	if err != nil {
		http.Error(w, "ошибка при обработке шаблона", http.StatusInternalServerError)
		return
	}

	t.Execute(w, s.db.All())
}
