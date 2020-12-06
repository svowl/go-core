package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
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
	g.spider = spider.New(10)
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

	return nil
}

// build запускает сканирование сайтов, сохранение данных в БД и построение индекса
func (g *gosearch) build() {
	if data, err := g.spider.Scan(g.urls, 2); err == nil {
		for _, doc := range data {
			// Добавление документа в БД. При этом для документа генерируется ID
			err := g.db.AddDoc(&doc)
			if err != nil {
				// Ошибка при добавлении документа - игнорируем его и идем дальше
				fmt.Printf("[warning] %v\n", err)
				continue
			}
			// Добавление документа в индекс
			g.index.Add(doc)
		}
	}
	// Проверяем есть ли вообще данные в БД и если нет - завершаем работу с ошибкой, т.к. это фатальная ситуация
	cnt := g.db.Count()
	if cnt == 0 {
		log.Fatal("[build] В БД нет документов!")
	}
	fmt.Printf("\n[build] В БД %d документов\n", cnt)
}

// run запускает веб-сервер.
func (g *gosearch) run() {
	g.router.HandleFunc("/", mainHandler).Methods(http.MethodGet)
	log.Println("Запуск http-сервера на интерфейсе:", ":80")
	srv := &http.Server{
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 40 * time.Second,
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

/*
// startServer запускает tcp сервер и обработчик клиентских подключений
func (g *gosearch) serve() {
	// регистрация сетевой службы
	listener, err := net.Listen("tcp4", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Сервер запущен на :8000")

	// цикл обработки клиентских подключений
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go g.handler(conn)
	}
}

// hander обрабатывает клиентские подключения
// зеркалим выводом в консоль для отладки
func (g *gosearch) handler(conn net.Conn) {
	defer conn.Close()
	defer fmt.Println("Connection closed")

	r := bufio.NewReader(conn)
	for {
		conn.SetDeadline(time.Now().Add(time.Second * 5))
		// Читаем поисковую фразу
		msg, _, err := r.ReadLine()
		if err != nil {
			return
		}
		phrase := string(msg)

		if phrase != "" {
			fmt.Println(phrase)
			// Выполняем поиск и пишем результат в соединение
			docs := g.engine.Search(phrase)
			resp := ""
			if docs == nil {
				resp = "No results\r\n"
			} else {
				for _, document := range docs {
					resp = resp + fmt.Sprintf("%s: %s\r\n", document.URL, document.Title)
				}
			}
			_, err = conn.Write([]byte(resp))
			if err != nil {
				return
			}
			fmt.Println(resp)
		}
	}
}
*/
