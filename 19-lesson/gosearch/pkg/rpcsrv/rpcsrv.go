package rpcsrv

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/engine"
)

// Service это структура RPC службы
type Service struct {
	engine *engine.Service
}

// New возвращает новый объект Service
func New(e *engine.Service) *Service {
	var s Service
	s.engine = e
	return &s
}

// Search вызывает поиск по query и возвращает список найденных документов в result
func (s *Service) Search(query string, result *[]crawler.Document) error {
	docs := s.engine.Search(query)
	*result = docs
	return nil
}
