package storage

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/storage/array"
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

// Документ для benchmark
var benchmarkDoc = crawler.Document{
	ID: 0, URL: "https://google.com", Title: "Google",
}

func TestDb_AddDocbTree(t *testing.T) {
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

func TestDb_AddDocArray(t *testing.T) {
	// Хранилище - память
	storage := mem.New()
	// Поисковая структура - бинарное дерево
	searcher := array.New()

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
	searcher2 := array.New()

	// Инициализируем с хранилищем в памяти и поисковой структурой - бинарным деревом
	db2.Init(storage, searcher2)

	// Проверяем поиск
	_, found = db2.Find(id)
	if found == false {
		t.Fatalf("Find(4): документ не найден (ID: %v)", id)
	}
}

func Benchmark_AddDocTree(b *testing.B) {
	storage := mem.New()
	searcher := btree.New()
	db := New()
	db.Init(storage, searcher)

	for i := 0; i < b.N; i++ {
		doc := benchmarkDoc
		err := db.AddDoc(&doc)
		if err != nil {
			b.Fatalf("Ошибка добавления документа: %v", err)
		}
	}
}

func Benchmark_AddDocArray(b *testing.B) {
	storage := mem.New()
	searcher := array.New()
	db := New()
	db.Init(storage, searcher)

	for i := 0; i < b.N; i++ {
		doc := benchmarkDoc
		err := db.AddDoc(&doc)
		if err != nil {
			b.Fatalf("Ошибка добавления документа: %v", err)
		}
	}
}

func Benchmark_FindTree(b *testing.B) {
	storage := mem.New()
	searcher := btree.New()
	db := New()
	db.Init(storage, searcher)
	var lastID int
	for i := 0; i < 10000; i++ {
		doc := benchmarkDoc
		err := db.AddDoc(&doc)
		if err != nil {
			b.Fatalf("Ошибка добавления документа: %v", err)
		}
		lastID = doc.ID
	}

	for i := 0; i < b.N; i++ {
		doc, found := db.Find(lastID)
		if found == false {
			b.Fatalf("Find(1): документ не найден (ID: %v)", lastID)
		}
		_ = doc
	}
}

func Benchmark_FindArray(b *testing.B) {
	storage := mem.New()
	searcher := array.New()
	db := New()
	db.Init(storage, searcher)
	var lastID int
	for i := 0; i < 10000; i++ {
		doc := benchmarkDoc
		err := db.AddDoc(&doc)
		if err != nil {
			b.Fatalf("Ошибка добавления документа: %v", err)
		}
		lastID = doc.ID
	}

	for i := 0; i < b.N; i++ {
		doc, found := db.Find(lastID)
		if found == false {
			b.Fatalf("Find(1): документ не найден (ID: %v)", lastID)
		}
		_ = doc
	}
}
