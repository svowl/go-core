package index

import (
	"gosearch/pkg/crawler"
	"reflect"
	"sort"
	"testing"
)

// Входные данные для тестирования
var testDocs = []crawler.Document{
	{ID: 1, URL: "https://google.com", Title: "Google"},
	{ID: 2, URL: "https://go.dev/", Title: "Why Go"},
	{ID: 3, URL: "https://golang.org/", Title: "The Go Programming Language"},
}

func TestIndex_Add(t *testing.T) {
	index := New()
	// Добавляем документы в индекс
	for _, doc := range testDocs {
		index.Add(doc)
	}
	// Проверяем - в индексе должен быть ключ "google", содержащий ссылку на документ с ID = 1
	if !exists(index.hash["google"], 1) {
		t.Errorf("Документ %v не добавлен", 1)
	}
	// В индексе ожидается 6 записей
	got := len(index.hash)
	exp := 6
	if got != exp {
		t.Log(index.hash)
		t.Errorf("Len(hash): получено %v, ожидается %v", got, exp)
	}

	// Добавляем документы повторно
	for _, doc := range testDocs {
		index.Add(doc)
	}
	// В индексе снова должно быть 6 слов, дублей быть не должно
	got = len(index.hash)
	if got != exp {
		t.Log(index.hash)
		t.Errorf("Len(hash): получено %v, ожидается %v", got, exp)
	}

	// Проверка на наличие дубликатов IDs
	for word, ids := range index.hash {
		if len(ids) < 2 {
			continue
		}
		// Сортируем список ID документов
		sort.Ints(ids)
		for i := 1; i < len(ids); i++ {
			// Ищем два последовательно одинаковых ID, если нашли - кричим ошибку
			if ids[i] == ids[i-1] {
				t.Errorf("Найден дубликат в списке ID (word: %v, ID: %v)", word, ids[i])
			}
		}
	}
}

func TestIndex_Search(t *testing.T) {
	index := New()
	// Добавляем документы в индекс
	for _, doc := range testDocs {
		index.Add(doc)
	}
	// Тест успешного поиска
	word := "go"
	gotIds := index.Search(word)
	expIds := []int{2, 3}
	if !reflect.DeepEqual(gotIds, expIds) {
		t.Errorf("Слово %v не найдено в индексе", word)
	}
	// Тест неудачного поиска
	word = "noword"
	gotIds = index.Search(word)
	if gotIds != nil {
		t.Errorf("Слово %v найдено в индексе", word)
	}
}
