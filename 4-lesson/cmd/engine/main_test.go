package main

import (
	"reflect"
	"testing"

	spider "go-core/4-lesson/pkg/fakespider"
	"go-core/4-lesson/pkg/index"
)

func Test_search(t *testing.T) {
	type args struct {
		phrase string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "first test",
			args: args{phrase: "google"},
			want: map[string]string{
				"https://www.google.com": "Google",
			},
		},
		{
			name: "second test",
			args: args{phrase: "Google"},
			want: map[string]string{
				"https://www.google.com": "Google",
			},
		},
		{
			name: "third test",
			args: args{phrase: "why"},
			want: map[string]string{
				"https://go.dev/solutions":              "Why Go - go.dev",
				"https://go.dev/solutions#case-studies": "Why Go - go.dev",
				"https://go.dev/solutions#use-cases":    "Why Go - go.dev",
			},
		},
		{
			name: "fourth test",
			args: args{phrase: "something not existing"},
			want: map[string]string{},
		},
	}

	// тестируем на фейковых данных
	s := new(spider.Spider)
	data, err := Scan(s, "", 2)
	if err != nil {
		t.Error("ошибка сканирования")
	}
	index.Build(data)
	index.SortRecords()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := search(tt.args.phrase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("получено %v, ожидается %v", got, tt.want)
			}
		})
	}
}
