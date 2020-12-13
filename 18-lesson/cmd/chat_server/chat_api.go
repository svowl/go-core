package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

// Service это служба чата.
// - chMessages: ассоциативный массив каналов для записи сообщений во все открытые соединения по /messages
// - nextConnID: ID следущего соединения
// - logger:     интерфейс для записи логов
type Service struct {
	upgrader   websocket.Upgrader
	chMessages map[int]chan string
	nextConnID int
	logger     io.Writer
}

// New возвращает новый объект службы
func New(logger io.Writer) *Service {
	var s Service
	s.upgrader = websocket.Upgrader{}
	s.chMessages = make(map[int]chan string)
	s.nextConnID = 0
	s.logger = logger
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
		for _, ch := range s.chMessages {
			ch <- string(message)
		}

		return
	}
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
	connID := s.nextConnID
	s.nextConnID++
	s.chMessages[connID] = make(chan string)

	// Пишем в канал текущего соединения
	for message := range s.chMessages[connID] {
		conn.WriteMessage(websocket.TextMessage, []byte(message))
		s.log(fmt.Sprintf("/messages: Write[conn %d]: %s", connID, message))
	}

	// Закрываем канал и удаляем его из s.chMessages
	close(s.chMessages[connID])
	delete(s.chMessages, connID)
}

// Логгер
func (s *Service) log(message string) {
	s.logger.Write([]byte(message + "\n"))
}
