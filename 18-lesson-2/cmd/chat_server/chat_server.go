package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	s := New(os.Stdout)
	s.endpoints()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write([]byte("Сервер запущен на localhost:8080"))
}
