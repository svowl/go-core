package membot

import (
	"gosearch/pkg/crawler"
)

// Service - имитация служба поискового робота.
type Service struct{}

// New - констрктор имитации службы поискового робота.
func New() *Service {
	s := Service{}
	return &s
}

// Scan возвращает заранее подготовленный набор данных
func (s *Service) Scan(url string, depth int) ([]crawler.Document, error) {

	data := []crawler.Document{
		{
			URL:   "https://yandex.ru",
			Title: "Яндекс",
		},
		{
			URL:   "https://www.google.ru",
			Title: "Google",
		},
		{
			URL:   "https://go.dev",
			Title: "go.dev",
		},
		{
			URL:   "https://go.dev/about",
			Title: "About - go.dev",
		},
		{
			URL:   "https://go.dev/learn",
			Title: "Learn - go.dev",
		},
		{
			URL:   "https://go.dev/solutions",
			Title: "Why Go - go.dev",
		},
	}

	return data, nil
}
