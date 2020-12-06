package storage

import (
	"gosearch/pkg/crawler"
)

// IReadWriter это интерфейс, объявляющий контракты чтения/записи данных в хранилище (файл или память)
type IReadWriter interface {
	Read() ([]byte, error)
	Write([]byte) (int, error)
}

// ISearcher это интерфейс, объявляющий контракты работы со структурой данных
type ISearcher interface {
	AddDoc(crawler.Document)
	FindDoc(int) (crawler.Document, bool)
	Count() int
	All() []crawler.Document
	JSONData() ([]byte, error)
	LoadFromJSON([]byte) error
	GenerateID() (int, error)
}

// Db - структура БД, содержит хранилище (file, mem) и поисковую структуру (btree)
type Db struct {
	storage  IReadWriter
	searcher ISearcher
}

// New создает объект Db
func New() *Db {
	var db Db
	return &db
}

// Init инициализирует объект Db переданными хранилищем и поисковой структурой
func (db *Db) Init(st IReadWriter, se ISearcher) {
	db.storage = st
	db.searcher = se

	json, err := db.storage.Read()
	if err == nil {
		// Если нет ошибки, загружаем данные из хранилища в поисковую структуру
		db.searcher.LoadFromJSON(json)
	}
}

// AddDoc добавляет документ в поисковую структуру
func (db *Db) AddDoc(doc *crawler.Document) error {
	if doc.ID == 0 {
		// Если ID не установлен, генерируем уникальный...
		ID, err := db.searcher.GenerateID()
		if err != nil {
			return err
		}
		doc.ID = ID
	}
	db.searcher.AddDoc(*doc)
	return nil
}

// Find выполняет поиск документа в поисковой структуре
func (db *Db) Find(id int) (crawler.Document, bool) {
	return db.searcher.FindDoc(id)
}

// Count возвращает кол-во документов в структуре
func (db *Db) Count() int {
	return db.searcher.Count()
}

// All возвращает все документы
func (db *Db) All() []crawler.Document {
	return db.searcher.All()
}

// Save сохраняет поисковую структуру в хранилище
func (db *Db) Save() (int, error) {
	jsonData, err := db.searcher.JSONData()
	if err != nil {
		return 0, err
	}
	return db.storage.Write(jsonData)
}

// Load загружает данные из хранилища в поисковую структуру
func (db *Db) Load() (int, error) {
	jsonData, err := db.storage.Read()
	if err != nil {
		return 0, err
	}
	db.searcher.LoadFromJSON(jsonData)
	return len(jsonData), nil
}
