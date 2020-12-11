package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var s *Service

func TestMain(m *testing.M) {
	s = New()
	s.router = mux.NewRouter()
	s.logger = os.Stdout
	s.endpoints()

	os.Exit(m.Run())
}

func Test_authJWT(t *testing.T) {
	// Тестирование пустого запроса, ожидаем ошибку 500
	req := httptest.NewRequest(http.MethodPost, "/auth", nil)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Status code: получено %d, ожидалось %d", rr.Code, http.StatusInternalServerError)
	}

	// Тестирование запроса с некорректными данными, ожидается ошибка 403
	data := authInfo{
		Login:    "john@example.com",
		Password: "111111",
	}
	payload, _ := json.Marshal(data)
	req = httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(payload))
	rr = httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Errorf("Status code: получено %d, ожидалось %d", rr.Code, http.StatusForbidden)
	}

	// Тестирование запроса с корректными данными
	data = authInfo{
		Login:    "john@example.com",
		Password: "12345",
	}
	payload, _ = json.Marshal(data)
	req = httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(payload))
	rr = httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Status code: получено %d, ожидалось %d", rr.Code, http.StatusOK)
	}
	// Декодирование JWT
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(rr.Body.String(), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretPassword), nil
	})
	if err != nil {
		t.Fatal("Ошибка парсинга токена", err)
	}

	// Проверка прав доступа, полученных из claims JWT токена
	var got []string
	err = json.Unmarshal([]byte(claims["perm"].(string)), &got)
	if err != nil {
		t.Fatalf("Ошибка декодирования прав доступа: %v", err)
		return
	}
	want := []string{"v", "c", "u", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Имя: получено %v, ожидается %v", got, want)
	}
}
