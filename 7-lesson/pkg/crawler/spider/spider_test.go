package spider

import (
	"testing"
)

func TestService_Scan(t *testing.T) {
	const url = "https://google.com/"
	s := New()
	data, err := s.Scan(url, 1)
	if err != nil {
		t.Fatal(err)
	}
	// Ожидаем в результатах только одну запись
	wantTitle := "Google"
	for _, v := range data {
		if v.URL != "https://google.com/" {
			t.Fatalf("получено %v, ожидается %v", v.URL, url)
		}
		if v.Title != wantTitle {
			t.Fatalf("получено %v, ожидается %v", v.Title, wantTitle)
		}
	}
	data, err = s.Scan(url, 2)
	if err != nil {
		t.Fatal(err)
	}
	// Ожидаем в результатах несколько записей с непустыми URL
	for _, v := range data {
		if v.URL == "" {
			t.Fatalf("Получен пустой URL")
		}
	}
}
