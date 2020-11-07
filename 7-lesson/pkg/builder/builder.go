package builder

import (
	"fmt"

	"go-core/7-lesson/pkg/index"
	"go-core/7-lesson/pkg/spider"
	"go-core/7-lesson/pkg/storage"
)

// Пакет builder запускает сканирование сайтов, передает полученные данные на индексирование,
// записывает индекс и данные в хранилище.

// Service хранит конфигурацию билдера: сайты для сканирования, сканер, хранилище, канал, индекс
// и флаг отключения вывода сообщений в терминал
type Service struct {
	Urls    []string
	Scanner spider.Scanner
	Storage storage.ReaderWriter
	channel chan<- int
	index   *index.Index
	silent  bool
}

// New создает объект Service
func New() *Service {
	s := Service{}
	return &s
}

// Init инициализирует объекст Service
func (s *Service) Init(urls []string, sp spider.Scanner, st storage.ReaderWriter) error {
	s.Urls = urls
	s.Scanner = sp
	s.Storage = st
	s.silent = false

	var err error
	// Создаем индекс для сканирования
	s.index, err = index.New(s.Storage)
	if err != nil {
		return err
	}

	return nil
}

// Channel инициализирует канал для сообщений о готовности данных
func (s *Service) Channel(ch chan<- int) {
	s.channel = ch
}

// Silent включает/выключает молчаливый режим работы пакета (см. функцию echo)
func (s *Service) Silent(flag bool) {
	s.silent = flag
}

// Build сканирует страницы, строит индекс
func (s *Service) Build() {
	for _, url := range s.Urls {
		s.echo("Сканируем  %s...", url)
		data, err := s.Scanner.Scan(url, 2)
		if err != nil {
			s.echo("ошибка сканирования сайта %s: %v", url, err)
			continue
		}
		s.echo("...найдено %d документов, индексируем...", len(data))
		// Строим индекс по списку просканированных документов
		_, err = s.index.Build(data)
		if err != nil {
			s.echo("%s", err)
			continue
		}
		// Сохраняем индекс в файл
		err = s.index.SaveData()
		if err != nil {
			s.echo("%s", err)
			continue
		}
		if s.channel != nil {
			// Посылаем сигнал в канал о необходимости обновить текущий поисковый индекс
			s.channel <- 1
		}
		s.echo("Проиндексировано %d страниц", s.index.Records.Count)
	}
}

// echo выводит сообщение в терминал, если silent mode выключен
func (s *Service) echo(format string, params ...interface{}) {
	if s.silent == false {
		fmt.Printf("\n[builder] "+format, params...)
	}
	// TODO: добавить запись в лог
}
