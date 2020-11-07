package main

import (
	"log"

	"go-core/7-lesson/pkg/builder"
	"go-core/7-lesson/pkg/engine"
	"go-core/7-lesson/pkg/spider"
	"go-core/7-lesson/pkg/storage"
	"go-core/7-lesson/pkg/storage/file"
)

// Пакет gosearch запускает два независимых процесса:
//  - сканирование + построение индекса (pkg.builder)
//  - получение пользовательского ввода и поиск (pkg.engine)
// Каждая часть работает со своим экземпляром индекса.
// Поисковый процесс обновляет свой индекс из файла после каждой записи туда новых данных билдером.

// gosearch хранит объекты builder и engine
type gosearch struct {
	builder *builder.Service
	engine  *engine.Service
	config  config
}

// config хранит конфигурацию для gosearch
type config struct {
	urls    []string
	storage storage.ReaderWriter
	scanner spider.Scanner
	channel chan int
}

func main() {
	// Создаем конфигурацию
	conf := config{
		urls:    []string{"https://www.google.com", "https://go.dev", "https://golang.org"},
		storage: file.New("./index.json"),
		scanner: spider.New(),
		channel: make(chan int, 1),
	}

	// Создаем объект gosearch
	g := new()
	// Инициализируем
	err := g.init(conf)
	if err != nil {
		log.Fatal(err)
	}

	g.run()
}

// Возвращает ссылку на новый объект gosearch
func new() *gosearch {
	g := gosearch{}
	g.builder = builder.New()
	g.engine = engine.New()
	return &g
}

// Инициализирует объект gosearch переданной конфигурацией
func (g *gosearch) init(conf config) error {
	g.config = conf

	// Инициализируем билдер, инициализиуем его urls, spider и storage
	err := g.builder.Init(conf.urls, conf.scanner, conf.storage)
	if err != nil {
		return err
	}

	if conf.channel != nil {
		// Инициализируем поисковик
		g.engine.Init(conf.storage)

		// Выключаем вывод сообщение от билдера
		g.builder.Silent(true)

		// Передаем канал для обмена сообщений в builder и engine
		g.builder.Channel(conf.channel)
		g.engine.Channel(conf.channel)
	}

	return nil
}

// Запускаем процесс...
func (g *gosearch) run() {

	if g.config.channel != nil {
		// С конфигурацией передан канал - запускаем builder и engine в отдельных параллельных потоках
		go g.builder.Build()
		// Запускаем слушателя в engine
		go g.engine.Listen()
		// Запускаем поисковый движок
		g.engine.Search()

	} else {
		// Иначе запускаем процессы последовательно: сначала builder, потом engine
		g.builder.Build()
		// Инициализируем поисковик, чтобы подхватить данные, созданные builder'ом
		g.engine.Init(g.config.storage)
		// Запускаем поисковый движок
		g.engine.Search()
	}
}
