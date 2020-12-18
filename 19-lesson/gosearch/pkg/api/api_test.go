package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gosearch/pkg/crawler"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/btree"
	"gosearch/pkg/storage/mem"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

// Объект webapp.Service
var s *Service

// Тестовые данные для БД и индекса
var docs = []crawler.Document{
	{ID: 1, URL: "https://google.com", Title: "Google"},
	{ID: 2, URL: "https://go.dev/", Title: "Why Go"},
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
	engine := engine.New()
	engine.Init(index, db)

	mux := mux.NewRouter()

	s = New(mux, db, index, engine)
	s.endpoints()

	os.Exit(m.Run())
}

func TestService_searchHandler(t *testing.T) {
	// Выполняем тестовый запрос поиска по фразе "google"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/google", nil)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("код неверен: получено %d, ожидалось %d", rr.Code, http.StatusOK)
	}
	// Декодируем
	var got []crawler.Document
	err := json.Unmarshal([]byte(rr.Body.String()), &got)
	if err != nil {
		t.Fatal("Ошибка декодирования данных", err)
	}
	// Проверяем ответ на соответствие ожиданиям
	want := []crawler.Document{
		{ID: 1, URL: "https://google.com", Title: "Google"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("получено %v, ожидается %v", got, want)
	}

	// Выполняем тестовый запрос поиска с пустыми результатами
	req = httptest.NewRequest(http.MethodGet, "/api/v1/search/not+existing+doc", nil)
	rr = httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("код неверен: получено %d, ожидалось %d", rr.Code, http.StatusOK)
	}
	// Декодируем
	err = json.Unmarshal([]byte(rr.Body.String()), &got)
	if err != nil {
		t.Fatal("Ошибка декодирования данных", err)
	}
	// Проверяем ответ на соответствие ожиданиям
	want = []crawler.Document{}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("получено %v, ожидается %v", got, want)
	}
}

func TestService_docsHandler(t *testing.T) {
	// Выполняем тестовый запрос
	req := httptest.NewRequest(http.MethodGet, "/api/v1/docs", nil)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("код неверен: получено %d, ожидалось %d", rr.Code, http.StatusOK)
	}
	// Декодируем
	var got []crawler.Document
	err := json.Unmarshal([]byte(rr.Body.String()), &got)
	if err != nil {
		t.Fatal("Ошибка декодирования данных", err)
	}
	// Проверяем ответ на соответствие ожиданиям
	want := docs
	// Сортировка нужна для точного сравнения с помощью reflect.DeepEqual(), т.к. ответ возвращает неотсортированные данные
	sort.Slice(got, func(i, j int) bool { return got[i].Title < got[j].Title })
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("получено %v, ожидается %v", got, want)
	}
}

func TestService_updateDocHandler(t *testing.T) {
	// Данные для обновления документа
	data := crawler.Document{
		URL:   "https://www.example.com",
		Title: "Example",
	}
	// Выполняем запрос с корректными данными, обновляется документ с ID=1
	payload, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/doc/1", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("код неверен: получено %d, ожидалось %d", rr.Code, http.StatusOK)
	}
	// Поиск документа в БД и сравнение с ожидаемым результатом
	id := 1
	got, found := s.db.Find(id)
	if found == false {
		t.Fatalf("Документ %v не найден", id)
	}
	want := crawler.Document{
		ID:    1,
		URL:   "https://www.example.com",
		Title: "Example",
		Body:  "",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Данные не верны: получено %v, ожидается %v", got, want)
	}

	// Запрос на обновление несуществующего документа, должен вернуть 404
	req = httptest.NewRequest(http.MethodPut, "/api/v1/doc/1000", bytes.NewBuffer(payload))
	rr = httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("код неверен: получено %d, ожидалось %d", rr.Code, http.StatusNotFound)
	}
}

func TestService_addDocHandler(t *testing.T) {
	// Данные для обновления документа
	data := crawler.Document{
		URL:   "https://www.example2.com",
		Title: "Example",
	}
	// Выполняем запрос
	payload, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/doc", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("код неверен: получено %d, ожидалось %d", rr.Code, http.StatusOK)
	}
	// Поиск документа в БД и сравнение с ожидаемым результатом
	found := false
	for _, doc := range s.db.All() {
		if doc.URL == data.URL && doc.Title == data.Title {
			found = true
			break
		}
	}
	if found == false {
		t.Fatal("Документ не создан")
	}
}

func TestService_deleteDocHandler(t *testing.T) {
	// Поиск документа в БД и сравнение с ожидаемым результатом
	id := 1
	_, found := s.db.Find(id)
	if found == false {
		t.Fatalf("Документ %v не найден", id)
	}
	// Выполняем запрос с на удаление документ с ID=1
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/doc/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получено %d, ожидалось %d", rr.Code, http.StatusOK)
	}
	// Поиск документа в БД и сравнение с ожидаемым результатом
	doc, found := s.db.Find(id)
	if found == true {
		t.Errorf("Документ %v найден, %v", id, doc)
	}

	// Запрос на обновление несуществующего документа, должен вернуть 404
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/doc/1000", nil)
	rr = httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("код неверен: получено %d, ожидалось %d", rr.Code, http.StatusNotFound)
	}
}
