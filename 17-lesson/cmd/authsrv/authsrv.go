package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Запускаем сервер
func main() {
	s := New()
	s.router = mux.NewRouter()
	s.logger = os.Stdout
	s.endpoints()
	http.ListenAndServe(":80", s.router)
}
