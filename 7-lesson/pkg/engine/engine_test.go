package engine

import (
	"errors"
	"reflect"
	"testing"

	"go-core/7-lesson/pkg/index"
	memstorage "go-core/7-lesson/pkg/storage/mem"
)

// Содержимое исходного файла с корректной json-строкой
var testContent = `{
	"Hash": {"go": [1, 2], "why": [2]},
	"Records": {
		"ID": 1,
		"Value": {"URL": "https://google.com", "Title": "Google"},
		"Left": null,
		"Right": {
			"ID": 2,
			"Value": {"URL": "https://go.dev", "Title": "Why Go"}
		}
	}
}`

func TestService_searchPhrase(t *testing.T) {
	e := New()

	// Инициализируем с ошибкой
	stErr := memstorage.New()
	stErr.Error = errors.New("Test error")
	err := e.Init(stErr)
	if err == nil {
		t.Fatalf("Ожидаем ошибку")
	}

	// Инициализируем без ошибки
	st := memstorage.New()
	st.Content = []byte(testContent)
	err = e.Init(st)
	if err != nil {
		t.Fatalf("Ошибка инициализации: %v", err)
	}

	tests := []struct {
		name string
		args string
		want []index.Record
	}{
		{
			name: "first test",
			args: "GO",
			want: []index.Record{
				{URL: "https://google.com", Title: "Google"},
				{URL: "https://go.dev", Title: "Why Go"},
			},
		},
		{
			name: "second test",
			args: "WHY",
			want: []index.Record{
				{URL: "https://go.dev", Title: "Why Go"},
			},
		},
		{
			name: "third test",
			args: "something not existing",
			want: []index.Record{},
		},
	}

	// проверяем результат
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.searchPhrase(tt.args); !reflect.DeepEqual(got, tt.want) {
				if tt.name == "third test" {
					if len(got) != 0 {
						t.Errorf("Фраза: %s, получено %v, ожидается %v", tt.args, reflect.ValueOf(got), reflect.ValueOf(tt.want))
					}
				} else {
					t.Errorf("Фраза: %s, получено %v, ожидается %v", tt.args, reflect.ValueOf(got), reflect.ValueOf(tt.want))
				}
			}
		})
	}
}
