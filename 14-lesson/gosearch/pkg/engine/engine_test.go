package engine

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"gosearch/pkg/crawler"
	"gosearch/pkg/index"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/btree"
	"gosearch/pkg/storage/mem"
)

// Входные данные для тестирования
var testDocs = []crawler.Document{
	{ID: 10, URL: "https://google.com", Title: "Google"},
	{ID: 12, URL: "https://go.dev/", Title: "Why Go"},
	{ID: 8, URL: "https://golang.org/", Title: "The Go Programming Language"},
	{ID: 18, URL: "https://yandex.ru/", Title: "Yandex"},
	{ID: 5, URL: "https://rambler.ru/", Title: "Rambler"},
	{ID: 11, URL: "https://mail.ru/", Title: "Mail.ru"},
}

var e *Service

func TestMain(m *testing.M) {
	// Хранилище - память
	mem := mem.New()
	// Поисковая структура - бинарное дерево
	btree := btree.New()
	// Создаем и инициализируем БД
	db := storage.New()
	db.Init(mem, btree)
	// Создаем индекс
	index := index.New()
	// Добавляем документы в БД и индекс
	for _, doc := range testDocs {
		err := db.AddDoc(&doc)
		if err != nil {
			fmt.Printf("Ошибка добавления документа в БД: %v", err)
			os.Exit(1)
		}
		index.Add(doc)
	}
	// Создаем и инициализируем объект engine
	e = New()
	e.Init(index, db)

	os.Exit(m.Run())
}

func TestService_Search(t *testing.T) {
	// Данные для тестирования
	tests := []struct {
		name   string
		phrase string
		want   []crawler.Document
	}{
		{
			name:   "Test #1",
			phrase: "Go",
			want: []crawler.Document{
				{ID: 12, URL: "https://go.dev/", Title: "Why Go"},
				{ID: 8, URL: "https://golang.org/", Title: "The Go Programming Language"},
			},
		},
		{
			name:   "Test #2",
			phrase: "WHY",
			want: []crawler.Document{
				{ID: 12, URL: "https://go.dev/", Title: "Why Go"},
			},
		},
	}

	// проверяем результат
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := e.Search(tt.phrase)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Фраза: %s, получено %v, ожидается %v", tt.phrase, reflect.ValueOf(got), reflect.ValueOf(tt.want))
			}
		})
	}

	phrase := "something not existing"
	if e.Search(phrase) != nil {
		t.Errorf("Фраза: %s найдена, хотя не должна", phrase)
	}
}
