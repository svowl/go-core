package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/gorilla/mux"

	"gosearch/pkg/api"
	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/rpcsrv"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/btree"
	"gosearch/pkg/storage/mem"
	"gosearch/pkg/webapp"
)

// gosearch содержит сканер, поисковый движок, БД, индекс, а также необходимые для их работы данные:
// хранилище (storage - file), поисковая структура (searcher - btree), список сайтов для сканирования
type gosearch struct {
	urls     []string
	spider   crawler.Scanner
	storage  storage.IReadWriter
	searcher storage.ISearcher
	db       *storage.Db
	index    *index.Index
	engine   *engine.Service
	router   *mux.Router
	webapp   *webapp.Service
	api      *api.Service
	rpcsrv   *rpcsrv.Service
}

func main() {

	// Сайты для сканирования
	var urls = []string{
		"https://google.com",
		"https://go.dev",
		"https://golang.org",
	}

	// Создание и инициализация gosearch
	// Для инициализации передается список сайтов и путь к файлу БД
	g := new()
	err := g.init(urls)
	if err != nil {
		log.Fatalf("Ошибка инициализации: %s", err)
	}

	// Запуск генерации поисковых данных
	g.build()

	// Запуск вебсервера
	g.run()
}

// new создает объект gosearch
func new() *gosearch {
	var g gosearch
	return &g
}

// init инициализирует объект gosearch
func (g *gosearch) init(urls []string) error {
	// Список сайтов
	g.urls = urls
	// Создаем объект сканер с максимум 10 потоками
	g.spider = spider.New()
	// Хранилище данных для БД (memory)
	g.storage = mem.New()
	// Поисковая структура для БД (бинарное дерево)
	g.searcher = btree.New()
	// БД
	g.db = storage.New()
	g.db.Init(g.storage, g.searcher)
	// Индекс
	g.index = index.New()
	// Поисковый движок
	g.engine = engine.New()
	g.engine.Init(g.index, g.db)

	g.router = mux.NewRouter()

	g.webapp = webapp.New(g.router, g.db, g.index)
	g.api = api.New(g.router, g.db, g.index, g.engine)

	g.rpcsrv = rpcsrv.New(g.engine)

	return nil
}

// build запускает сканирование сайтов, сохранение данных в БД и построение индекса
func (g *gosearch) build() {
	// Запускаем многопоточное сканирование
	chRes, chErr := g.spider.BatchScan(g.urls, 2, 10)
	// Принимаем отсканированные документы по одному и добавляем в БД и индекс
	go func() {
		for doc := range chRes {
			err := g.db.AddDoc(&doc)
			if err != nil {
				// Ошибка при добавлении документа - игнорируем и идем дальше
				log.Println("Ошибка добавления документа в БД:", err)
				continue
			}
			g.index.Add(doc)
		}
	}()
	go func() {
		for range chErr {
			// Игнориуем ошибки сканирования
		}
	}()
}

// run запускает веб-сервер.
func (g *gosearch) run() {

	// Добавляем обработчик RPC-службы
	rpcserver := rpc.NewServer()
	err := rpcserver.Register(g.rpcsrv)
	if err != nil {
		log.Fatal(err)
	}
	g.router.Handle("/rpc", rpcserver)

	// Добавляем обработчик дефолтного пути "/"
	g.router.HandleFunc("/", mainHandler).Methods(http.MethodGet)

	log.Println("Запуск http-сервера на интерфейсе:", ":80")
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      g.router,
		Addr:         ":80",
	}
	listener, err := net.Listen("tcp4", srv.Addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.Serve(listener))
}

// HTTP-обработчик по умолчанию
func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html><body><h2>Gosearch Web App</h2></body></html>`))
}
