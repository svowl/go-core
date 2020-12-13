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
	s.nextConnID = 0
	s.logger = logger
	s.chCloseConn = make(chan int)
	return &s
}

// enpoints объявляет конечные точки
func (s *Service) endpoints() {
	http.HandleFunc("/send", s.sendHandler)
	http.HandleFunc("/messages", s.messagesHandler)
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

	authorized := false

	for {
		// Читаем сообщение из соединения
		mt, message, err := conn.ReadMessage()
		if err != nil {
			conn.WriteMessage(mt, []byte(err.Error()))
			s.log("/send: Read " + err.Error())
			return
		}

		// Первое сообщение должно быть "password". Если это не так, пишем в соединение сообщение об ошибке
		if !authorized && string(message) != "password" {
			conn.WriteMessage(mt, []byte("not authorized"))
			s.log("/send: Write: not authorized")
			return
		}

		// Отвечаем на "password" сообщением "OK" и переходим к чтению след. сообщения
		if !authorized && string(message) == "password" {
			conn.WriteMessage(websocket.TextMessage, []byte("OK"))
			s.log("/send: Write: OK")
			authorized = true
			continue
		}

		s.log("/send: Message received: " + string(message))

		// Пишем в сообщение в каналы всех открытых соединений
		// Поскольку читаем из общей памяти (s.chMessages), ставим лок
		s.mux.Lock()
		for _, ch := range s.chMessages {
			ch <- string(message)
		}
		s.mux.Unlock()

		return
	}
}

// Обработчик для /messages пишет в соединение сообщения, принятые через /send
func (s *Service) messagesHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from ", r)
		}
	}()
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

	// Пишем в канал текущего соединения
	for message := range connChannel {
		conn.WriteMessage(websocket.TextMessage, []byte(message))
		s.log(fmt.Sprintf("/messages: Write[conn %d]: %s", connID, message))
	}

	// Во избежание паники закрытие канала и удаление из списка каналов вынесено в отдельный метод и синхронизируется
	// через отдельный канал
	s.chCloseConn <- connID
}

// Логгер
func (s *Service) log(message string) {
	s.logger.Write([]byte(message + "\n"))
}

// closeMsgChan закрывает канал сообщений соединения и удаляет его из s.chMessages
func (s *Service) closeMsgChan(ID int) {
	for connID := range s.chCloseConn {
		s.mux.Lock()
		close(s.chMessages[connID])
		delete(s.chMessages, connID)
		s.mux.Unlock()
	}
}
