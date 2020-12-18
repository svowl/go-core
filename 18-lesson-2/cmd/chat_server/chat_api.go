package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Service это служба чата.
// - chMessages: ассоциативный массив каналов для записи сообщений во все открытые соединения по /messages
// - nextConnID: ID следущего соединения
// - logger:     интерфейс для записи логов
// - mux:        мьютекс нужен для установки лока при изменении общей памяти
// - closeMsgChan: вспомогательный канал для очистки списка каналов соединений chMessages после закрытия соединения
type Service struct {
	upgrader    websocket.Upgrader
	chMessages  map[int]chan string
	nextConnID  int
	logger      io.Writer
	mux         sync.Mutex
	chCloseConn chan int
}

// New возвращает новый объект службы
func New(logger io.Writer) *Service {
	var s Service
	s.upgrader = websocket.Upgrader{}
	s.chMessages = make(map[int]chan string)
	s.logger = logger
	s.chCloseConn = make(chan int)
	return &s
}

// enpoints объявляет конечные точки
func (s *Service) endpoints() {
	http.HandleFunc("/send", s.sendHandler)
	http.HandleFunc("/messages", s.messagesHandler)
	go s.closeMsgChan()
}

// Обработчик для /send принимает сообщение от пользователя
func (s *Service) sendHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.log("/send: Ошибка соединения: " + err.Error())
		return
	}
	defer conn.Close()

	// Читаем сообщение из соединения
	mt, message, err := conn.ReadMessage()
	if err != nil {
		conn.WriteMessage(mt, []byte(err.Error()))
		s.log("/send: Read " + err.Error())
		return
	}

	// Первое сообещение ожидается "password". Отечаем на него сообщением "OK"
	if string(message) == "password" {
		conn.WriteMessage(websocket.TextMessage, []byte("OK"))
		s.log("/send: Write: OK")

	} else {
		conn.WriteMessage(mt, []byte("not authorized"))
		s.log("/send: Write: not authorized")
		return
	}

	// Читаем следующее сообщение
	mt, message, err = conn.ReadMessage()
	if err != nil {
		conn.WriteMessage(mt, []byte(err.Error()))
		s.log("/send: Read " + err.Error())
		return
	}

	s.log("/send: Message received: " + string(message))

	// Пишем в сообщение в каналы всех открытых соединений
	// Поскольку читаем из общей памяти (s.chMessages), ставим лок
	s.mux.Lock()
	for _, ch := range s.chMessages {
		ch <- string(message)
	}
	s.mux.Unlock()
}

// Обработчик для /messages пишет в соединение сообщения, принятые через /send
func (s *Service) messagesHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.log("/messages: Ошибка соединения: " + err.Error())
		return
	}
	defer conn.Close()

	// Увеличиваем счетчик соединений и создаем новый канал для записи
	s.mux.Lock()
	connID := s.nextConnID
	s.nextConnID++
	connChannel := make(chan string)
	s.chMessages[connID] = connChannel
	s.mux.Unlock()

	defer func() {
		// Во избежание паники закрытие канала и удаление из списка каналов вынесено в отдельный метод и синхронизируется
		// через канал s.chCloseConn
		s.log(fmt.Sprintf("Завершено соединение %d", connID))
		s.chCloseConn <- connID
	}()

	// Пишем в канал текущего соединения
	for message := range connChannel {
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			return
		}
		s.log(fmt.Sprintf("/messages: Write[conn %d]: %s", connID, message))
	}
}

// Логгер
func (s *Service) log(message string) {
	s.logger.Write([]byte(message + "\n"))
}

// closeMsgChan закрывает канал сообщений соединения и удаляет его из s.chMessages
func (s *Service) closeMsgChan() {
	for connID := range s.chCloseConn {
		s.mux.Lock()
		close(s.chMessages[connID])
		delete(s.chMessages, connID)
		s.mux.Unlock()
		s.log(fmt.Sprintf("Закрыт канал %d", connID))
	}
}
