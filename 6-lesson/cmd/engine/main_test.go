package main

import (
	"reflect"
	"testing"

	"go-core/6-lesson/pkg/index"
	"go-core/6-lesson/pkg/spider/memspider"
	"go-core/6-lesson/pkg/storage/mem"
)

func Test_search(t *testing.T) {
	type args struct {
		phrase string
	}
	tests := []struct {
		name string
		args args
		want []index.Record
	}{
		{
			name: "first test",
			args: args{phrase: "google"},
			want: []index.Record{
				{URL: "https://www.google.com", Title: "Google"},
			},
		},
		{
			name: "second test",
			args: args{phrase: "Google"},
			want: []index.Record{
				{URL: "https://www.google.com", Title: "Google"},
			},
		},
		{
			name: "third test",
			args: args{phrase: "something not existing"},
			want: []index.Record{},
		},
	}

	// тестируем на заранее подготовленных данных (пакет memspider)
	s := new(memspider.Spider)
	data, err := s.Scan("", 2)
	if err != nil {
		t.Error("ошибка сканирования")
	}
	// используем хранилище данных в памяти (storage/mem), изначально пустое
	storage := mem.Storage{}
	i, err := index.New(&storage)
	if err != nil {
		t.Error(err)
	}
	// строим индекс
	_, err = i.Build(data)
	if err != nil {
		t.Fatalf("Ошибка вызова Build(): %v", err)
	}
	err = i.SaveData()
	if err != nil {
		t.Fatalf("Ошибка SaveData(): %v", err)
	}
	// читаем индекс и разворачиваем в памяти текущий индекс для поиска
	updateCurrentIndex(&storage)
	// проверяем результат
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchPhrase(tt.args.phrase); !reflect.DeepEqual(got, tt.want) {
				if tt.name == "third test" {
					if len(got) != 0 {
						t.Errorf("получено %v, ожидается %v", reflect.ValueOf(got), reflect.ValueOf(tt.want))
					}
				} else {
					t.Errorf("получено %v, ожидается %v", reflect.ValueOf(got), reflect.ValueOf(tt.want))
				}
			}
		})
	}
}
