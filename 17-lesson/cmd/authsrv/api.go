package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// JWT secret password.
// По-хорошему, его надо хранить в .env
const secretPassword = "98thv23@49v3-"

// Service это служба авторизации
type Service struct {
	router *mux.Router
	logger io.Writer
}

// Структура, описывающая пользователя
type user struct {
	name, login, password string
	permissions           []string
}

// Структура принимающая данные авторизации из тела запроса
type authInfo struct {
	Login    string
	Password string
}

// Список пользователей содержит информацию об имени, логине (email), пароле и правах доступа.
// Права могут быть реализованы различными способами, здесь просто массив строковых идентификаторов,
// где v - view, c - create, u - update, d - delete.
var users = []user{
	{
		name:        "John",
		login:       "john@example.com",
		password:    "12345",
		permissions: []string{"v", "c", "u", "d"},
	},
	{
		name:        "Lisa",
		login:       "lisa@example.com",
		password:    "54321",
		permissions: []string{"v"},
	},
}

// New создает объект службы авторизации
func New() *Service {
	var s Service
	return &s
}

// endpoints объявляет логгер и конечную точку для авторизации
func (s *Service) endpoints() {
	s.router.Use(s.logMiddleware)
	s.router.HandleFunc("/auth", s.authJWT).Methods(http.MethodPost)
}

// authJWT проверяет данные запроса к /auth, генерирует JWT и возвращает его в ответе
func (s *Service) authJWT(w http.ResponseWriter, r *http.Request) {
	// Получаем данные из запроса
	var auth authInfo
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ищем пользователя по данным авторизации
	var authUser user
	for _, user := range users {
		if auth.Login == user.login && auth.Password == user.password {
			authUser = user
			break
		}
	}
	if false || authUser.login == "" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Подготовка данных о правах доступа
	perms, err := json.Marshal(authUser.permissions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  authUser.login,
		"name": authUser.name,
		"perm": string(perms),
		"nbf":  time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secretPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}
