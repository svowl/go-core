package engine

import (
	"reflect"
	"testing"

	"go-core/7-lesson/pkg/crawler"
	"go-core/7-lesson/pkg/index"
	memstorage "go-core/7-lesson/pkg/storage/mem"
)

// Содержимое исходного файла с корректной json-строкой
var testContent = `{
	"Hash": {"go": [1, 2], "why": [2, 3]},
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

func TestServiceSearch(t *testing.T) {
	e := New()

	st := memstorage.New()
	st.Content = []byte(testContent)

	ind := index.New()
	err := ind.Init(st)
	if err != nil {
		t.Fatalf("Ошибка инициализации индекса: %s", err)
	}

	e.Init(ind)

	tests := []struct {
		name string
		args string
		want []crawler.Document
	}{
		{
			name: "first test",
			args: "GO",
			want: []crawler.Document{
				{URL: "https://google.com", Title: "Google"},
				{URL: "https://go.dev", Title: "Why Go"},
			},
		},
		{
			name: "second test",
			args: "WHY",
			want: []crawler.Document{
				{URL: "https://go.dev", Title: "Why Go"},
			},
		},
		{
			name: "third test",
			args: "something not existing",
			want: []crawler.Document{},
		},
	}

	// проверяем результат
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, found := e.Search(tt.args); !found || !reflect.DeepEqual(got, tt.want) {
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
