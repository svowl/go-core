package engine

import (
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

func TestServiceSearch(t *testing.T) {

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
			t.Errorf("Ошибка добавления документа в БД: %v", err)
		}
		index.Add(doc)
	}

	// Создаем объект engine
	e := New()
	e.Init(index, db)

	// Данные для тестирования
	tests := []struct {
		name string
		args string
		want []crawler.Document
	}{
		{
			name: "Test #1",
			args: "Go",
			want: []crawler.Document{
				{ID: 12, URL: "https://go.dev/", Title: "Why Go"},
				{ID: 8, URL: "https://golang.org/", Title: "The Go Programming Language"},
			},
		},
		{
			name: "Test #2",
			args: "WHY",
			want: []crawler.Document{
				{ID: 12, URL: "https://go.dev/", Title: "Why Go"},
			},
		},
		{
			name: "Test #3",
			args: "something not existing",
			want: nil,
		},
	}

	// проверяем результат
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, found := e.Search(tt.args); !found || !reflect.DeepEqual(got, tt.want) {
				if tt.name == "Test #3" {
					if found {
						t.Errorf("Фраза: %s, получено %v, ожидается %v", tt.args, reflect.ValueOf(got), reflect.ValueOf(tt.want))
					}
				} else {
					t.Errorf("Фраза: %s, получено %v, ожидается %v", tt.args, reflect.ValueOf(got), reflect.ValueOf(tt.want))
				}
			}
		})
	}
}
