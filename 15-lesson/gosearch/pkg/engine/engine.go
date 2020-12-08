package engine

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/index"
	"gosearch/pkg/storage"
)

// Пакет engine производит поиск строки в индексе и возвращает список документов

// Service это структура для хранения состояния поискового движка
type Service struct {
	index   *index.Index
	storage *storage.Db
}

// New создает объект Service
func New() *Service {
	s := Service{}
	return &s
}

// Init инициализирует объект Service
func (s *Service) Init(index *index.Index, storage *storage.Db) {
	s.index = index
	s.storage = storage
}

// Search ищет слово в индексе и возвращает список документов
func (s *Service) Search(word string) []crawler.Document {
	var result []crawler.Document
	ids := s.index.Search(word)
	if ids == nil {
		// Слово не найдено, возвращаем false
		return nil
	}
	// Слово найдено в хеше, ids содержит индексы документов в БД
	for _, id := range ids {
		// Поиск записей в БД
		record, ok := s.storage.Find(id)
		if ok == false {
			// Не нашли документ в БД, продолжаем
			continue
		}
		// Документ найден - добавляем его в результат
		result = append(result, record)
	}
	return result
}
