package webapp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"gosearch/pkg/crawler"
	"gosearch/pkg/index"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/btree"
	"gosearch/pkg/storage/mem"
)

// Объект webapp.Service
var s *Service

// Тестовые данные для БД и индекса
var docs = []crawler.Document{
	{ID: 10, URL: "https://google.com", Title: "Google"},
	{ID: 12, URL: "https://go.dev/", Title: "Why Go"},
	{ID: 8, URL: "https://golang.org/", Title: "The Go Programming Language"},
	{ID: 18, URL: "https://yandex.ru/", Title: "Yandex"},
	{ID: 5, URL: "https://rambler.ru/", Title: "Rambler"},
	{ID: 11, URL: "https://mail.ru/", Title: "Mail.ru"},
}

func TestMain(m *testing.M) {
	db := storage.New()
	db.Init(mem.New(), btree.New())
	index := index.New()
	// Добавляем документы в БД и индекс
	for _, doc := range docs {
		err := db.AddDoc(&doc)
		if err != nil {
			fmt.Printf("Ошибка добавления документа в БД: %v", err)
			os.Exit(1)
		}
		index.Add(doc)
	}

	mux := mux.NewRouter()

	s = New(mux, db, index)
	s.endpoints()

	os.Exit(m.Run())
}

func TestService_indexHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/index", nil)
	rr := httptest.NewRecorder()

	s.router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	want := fmt.Sprintf("<h3>Total rows: %d</h3>", len(s.index.Hash))
	if false == strings.Contains(rr.Body.String(), want) {
		t.Errorf("Строка %v не найдена", want)
	}
}

func TestService_docsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rr := httptest.NewRecorder()

	s.router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	want := fmt.Sprintf("<h3>Total documents: %d</h3>", s.db.Count())
	if false == strings.Contains(rr.Body.String(), want) {
		t.Errorf("Строка %v не найдена", want)
	}
}
