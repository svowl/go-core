package array

import (
	"encoding/json"

	"gosearch/pkg/crawler"
)

// Data хранит массив документов
type Data struct {
	Docs []crawler.Document
}

// New создает объект array.Data
func New() *Data {
	var d Data
	d.Docs = make([]crawler.Document, 0, 10)
	return &d
}

// AddDoc добавляет документ в массив
func (d *Data) AddDoc(doc crawler.Document) {
	d.Docs = append(d.Docs, doc)
}

// FindDoc ищет документ и возвращает документ и флаг (true, если найден).
func (d *Data) FindDoc(id int) (crawler.Document, bool) {
	for i := 0; i < len(d.Docs); i++ {
		if d.Docs[i].ID == id {
			return d.Docs[i], true
		}
	}
	return crawler.Document{}, false
}

// Count возвращает кол-во документов.
func (d *Data) Count() int {
	//fmt.Println(d.Docs)
	return len(d.Docs)
}

// JSONData возвращает структуру array.Data в сериализованном виде
func (d *Data) JSONData() ([]byte, error) {
	jsonData, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// LoadFromJSON загружает в array.Data данные из json-строки
func (d *Data) LoadFromJSON(data []byte) error {
	err := json.Unmarshal(data, d)
	if err != nil {
		return err
	}
	return nil
}

// GenerateID генерирует уникальный ID документа.
func (d *Data) GenerateID() (int, error) {
	var ID int = 1
	for i := 0; i < len(d.Docs); i++ {
		if ID < d.Docs[i].ID {
			ID = d.Docs[i].ID + 1
		}
	}
	return ID, nil
}
