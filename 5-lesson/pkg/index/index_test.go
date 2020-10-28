package index

import (
	"reflect"
	"testing"
)

// Входные данные от spider
var inputData1 = map[string]string{
	"https://go.dev":                        "go.dev",
	"https://go.dev/":                       "go.dev",
	"https://go.dev/about":                  "About - go.dev",
	"https://go.dev/learn":                  "Learn - go.dev",
	"https://go.dev/solutions":              "Why Go - go.dev",
	"https://go.dev/solutions#case-studies": "Why Go - go.dev",
	"https://go.dev/solutions#use-cases":    "Why Go - go.dev",
	"https://www.google.com":                "Google",
}

// Входные данные от spider (второй набор) для тестирования уникальности URL'ов
var inputData2 = map[string]string{
	"https://go.dev":  "go.dev",
	"https://go.dev/": "go.dev",
}

func TestSearch(t *testing.T) {
	type args struct {
		phrase string
	}
	tests := []struct {
		name string
		args args
		want []Record
	}{
		{
			name: "first test",
			args: args{phrase: "google"},
			want: []Record{
				{
					URL:   "https://www.google.com",
					Title: "Google",
				},
			},
		},
		{
			name: "second test",
			args: args{phrase: "Google"},
			want: []Record{
				{
					URL:   "https://www.google.com",
					Title: "Google",
				},
			},
		},
		{
			name: "third test",
			args: args{phrase: "why"},
			want: []Record{
				{
					URL:   "https://go.dev/solutions",
					Title: "Why Go - go.dev",
				},
				{
					URL:   "https://go.dev/solutions#case-studies",
					Title: "Why Go - go.dev",
				},
				{
					URL:   "https://go.dev/solutions#use-cases",
					Title: "Why Go - go.dev",
				},
			},
		},
	}
	Build(inputData1)
	Build(inputData2)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Search(tt.args.phrase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search(%s) = %v, want %v", tt.args.phrase, reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
	got := Search("something not existing")
	if len(got) > 0 {
		t.Errorf("Search(something not existing) вернул неправильный результат %v", got)
	}
}
