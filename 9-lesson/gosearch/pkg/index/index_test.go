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

func TestIndex_Index(t *testing.T) {
	index := New()
	// Добавляем документы в индекс
	for _, doc := range testDocs {
		index.Add(doc)
		if !exists(index.hash["google"], 1) {
			t.Errorf("Add(%v): документ не добавлен", doc.ID)
		}
	}
	// В индексе должно быть 6 слов
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
	// В индексе снова должно быть 6 слов
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
		sort.Ints(ids)
		for i := 1; i < len(ids); i++ {
			if ids[i] == ids[i-1] {
				t.Errorf("Найден дубликат в списке ID (word: %v, ID: %v)", word, ids[i])
			}
		}
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
