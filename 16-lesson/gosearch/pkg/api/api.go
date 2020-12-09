package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"gosearch/pkg/crawler"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/storage"
)

// Service это служба Web-приложения, содержит ссылки на объекты роутера, БД и индекса
type Service struct {
	router *mux.Router
	db     *storage.Db
	index  *index.Index
	engine *engine.Service
}

// New создает объект Service, объявляет endpoints
func New(router *mux.Router, db *storage.Db, index *index.Index, engine *engine.Service) *Service {
	var s Service
	s.router = router
	s.db = db
	s.index = index
	s.engine = engine

	s.endpoints()

	return &s
}

// Определяем endpoints
func (s *Service) endpoints() {
	r := s.router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	r.HandleFunc("/search/{phrase}", s.searchHandler).Methods(http.MethodGet)
	r.HandleFunc("/search", s.searchHandler).Methods(http.MethodGet)
	r.HandleFunc("/docs", s.docsHandler).Methods(http.MethodGet)
	r.HandleFunc("/doc/{id:[0-9]+}", s.updateDocHandler).Methods(http.MethodPut)
	r.HandleFunc("/doc/{id:[0-9]+}", s.deleteDocHandler).Methods(http.MethodDelete)
	r.HandleFunc("/doc", s.addDocHandler).Methods(http.MethodPost)
}

// HTTP-обработчик api/v1/search/{phrase} выводит JSON-кодированные результаты поиска по phrase
func (s *Service) searchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docs := s.engine.Search(vars["phrase"])
	var err error
	if len(docs) > 0 {
		err = json.NewEncoder(w).Encode(docs)
	} else {
		err = json.NewEncoder(w).Encode([]string{})
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HTTP-обработчик /api/v1/docs выводит JSON-закодированный список документов в БД
func (s *Service) docsHandler(w http.ResponseWriter, r *http.Request) {
	docs := s.db.All()
	err := json.NewEncoder(w).Encode(docs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HTTP-обработчик PUT:/api/v1/doc/{id} обновляет документ в БД
func (s *Service) updateDocHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID документа из query запроса
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Проверяем существование документа по ID. Сам документ нам не нужен, только факт наличия
	_, found := s.db.Find(id)
	if found == false {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}
	// Получаем данные для обновления из тела запроса
	var data crawler.Document
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.ID = id
	// Обновляем (метод AddDoc() обновляет документ в БД, если ID уже существует)
	err = s.db.AddDoc(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("ok")
}

// HTTP-обработчик DELETE:/api/v1/doc/{id} удаляет документ из БД
func (s *Service) deleteDocHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID документа из query запроса
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Удаляем документ. Если не найден, возвращаем 404
	if s.db.DeleteDoc(id) == false {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode("ok")
}

// HTTP-обработчик POST:/api/v1/doc добавляет документ в БД
func (s *Service) addDocHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем данные для обновления из тела запроса
	var doc crawler.Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	doc.ID = 0
	// Добавляем документ (ID сгенерится при добавлении)
	err = s.db.AddDoc(&doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.index.Add(doc)
	json.NewEncoder(w).Encode(map[string]int{"id": doc.ID})
}
