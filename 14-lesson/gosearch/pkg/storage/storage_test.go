package storage

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/storage/btree"
	"gosearch/pkg/storage/mem"
	"testing"
)

// Входные данные для тестирования
var testDocs = []crawler.Document{
	{ID: 10, URL: "https://google.com", Title: "Google"},
	{ID: 12, URL: "https://go.dev/", Title: "Why Go"},
	{ID: 8, URL: "https://golang.org/", Title: "The Go Programming Language"},
	{ID: 18, URL: "https://yandex.ru/", Title: "Yandex"},
	{ID: 5, URL: "https://rambler.ru/", Title: "Rambler"},
	{ID: 0, URL: "https://mail.ru/", Title: "Mail.ru"},
}

func TestDb_AddDoc(t *testing.T) {
	// Хранилище - память
	storage := mem.New()
	// Поисковая структура - бинарное дерево
	searcher := btree.New()

	// Создаем и инициализируем объект БД
	db := New()
	db.Init(storage, searcher)

	// Добавляем документы
	for _, doc := range testDocs {
		err := db.AddDoc(&doc)
		if err != nil {
			t.Errorf("Ошибка добавления документа: %v", err)
		}
	}

	// Проверяем кол-во элементов
	if db.Count() != 6 {
		t.Fatalf("Count: получено %v, ожидается %v", db.Count(), 6)
	}

	// Проверяем поиск
	id := 18
	_, found := db.Find(id)
	if found == false {
		t.Fatalf("Find(1): документ не найден (ID: %v)", id)
	}

	// Сохраняем данные в хранилище
	_, err := db.Save()
	if err != nil {
		t.Fatalf("Ошибка при сохранении данных: %v", err)
	}

	// Загружаем из хранилища
	_, err = db.Load()
	if err != nil {
		t.Fatalf("Ошибка при загрузке данных: %v", err)
	}

	// Проверяем поиск
	_, found = db.Find(id)
	if found == false {
		t.Fatalf("Find(2): документ не найден (ID: %v)", id)
	}

	// Создаем второй объект БД
	db2 := New()
	searcher2 := btree.New()

	// Инициализируем с хранилищем в памяти и поисковой структурой - бинарным деревом
	db2.Init(storage, searcher2)

	// Проверяем поиск
	_, found = db2.Find(id)
	if found == false {
		t.Fatalf("Find(4): документ не найден (ID: %v)", id)
	}
}
