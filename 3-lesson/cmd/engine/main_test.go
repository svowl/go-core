package main

import (
	"reflect"
	"testing"

	spider "go-core/3-lesson/pkg/fakespider"
)

func Test_search(t *testing.T) {
	// тестируем на фейковых данных
	s := new(spider.Spider)
	data, err := Scan(s, "", 2)
	if err != nil {
		t.Error("ошибка сканирования")
	}
	type args struct {
		phrase  string
		storage map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "first test",
			args: args{phrase: "google", storage: data},
			want: map[string]string{
				"https://www.google.com": "Google",
			},
		},
		{
			name: "second test",
			args: args{phrase: "Google", storage: data},
			want: map[string]string{
				"https://www.google.com": "Google",
			},
		},
		{
			name: "third test",
			args: args{phrase: "why", storage: data},
			want: map[string]string{
				"https://go.dev/solutions":              "Why Go - go.dev",
				"https://go.dev/solutions#case-studies": "Why Go - go.dev",
				"https://go.dev/solutions#use-cases":    "Why Go - go.dev",
			},
		},
		{
			name: "fourth test",
			args: args{phrase: "something not existing", storage: data},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := search(tt.args.phrase, tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("получено %v, ожидается %v", got, tt.want)
			}
		})
	}
}
