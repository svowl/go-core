package index

import (
	"strings"

	"gosearch/pkg/crawler"
)

// Index это структура для хранения состояния индекса
type Index struct {
	hash map[string][]int
}

// New создает объект Index
func New() *Index {
	var index Index
	index.hash = make(map[string][]int)
	return &index
}

// Add добавляет документ в индекс
func (index *Index) Add(doc crawler.Document) {
	for _, word := range words(doc.Title) {
		if _, found := index.hash[word]; found && exists(index.hash[word], doc.ID) {
			continue
		}
		index.hash[word] = append(index.hash[word], doc.ID)
	}
}

// Search ищет по слову и возвращает список ID соответствующих документов в индексе
func (index *Index) Search(word string) []int {
	return index.hash[strings.ToLower(word)]
}

// words разделяет text на слова и возвращает в виде массива строк
func words(text string) []string {
	words := make([]string, 0, 10)
	for _, word := range strings.Fields(strings.ToLower(text)) {
		words = append(words, word)
	}
	return words
}

// exists возвращает true, если в массиве ids найдено число id
func exists(ids []int, id int) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}
