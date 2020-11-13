package index

import (
	"errors"
	"go-core/7-lesson/pkg/crawler"
	"go-core/7-lesson/pkg/storage/mem"
	"reflect"
	"testing"
)

// Содержимое исходного файла с некорректной json-строкой
var testContentWrong = "Wrong json data"

// Содержимое исходного файла с корректной json-строкой
var testContentCorrect = `{
	"Hash": {"go": [1, 2]},
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

// Входные данные для Build
var testInputData = []crawler.Document{
	{URL: "https://google.com", Title: "Google"},
	{URL: "https://go.dev/", Title: "Why Go"},
}

func TestNew(t *testing.T) {
	// Empty storage
	storage := mem.Storage{}
	ind := New()
	err := ind.Init(&storage)
	if err != nil {
		t.Error(err)
	}
	// Storage с некорректными json-данными
	storage = mem.Storage{
		Content: []byte(testContentWrong),
	}
	err = ind.Init(&storage)
	if err == nil {
		t.Errorf("Ожидается ошибка json, получено nil")
	}
	// Storage с ошибкой
	testError := errors.New("Test error")
	storage = mem.Storage{
		Error: testError,
	}
	err = ind.Init(&storage)
	if err != testError {
		t.Errorf("Ожидается ошибка %v, получено %v", testError, err)
	}
	// Storage с корректной json-строкой
	storage = mem.Storage{
		Content: []byte(testContentCorrect),
	}
	err = ind.Init(&storage)
	if err != nil {
		t.Errorf("Ожидается: нет ошибки, получено: %v", err)
	}
	want := []int{1, 2}
	if got, found := ind.Search("go"); found && !reflect.DeepEqual(got, want) {
		t.Errorf("Search: получено %v, ожидается %v", got, want)
	}
}

func TestIndex_Build(t *testing.T) {
	// Создаем индекс с пустым исходным файлом
	storage := mem.Storage{}
	ind := New()
	err := ind.Init(&storage)
	if err != nil {
		t.Error(err)
	}
	// Строим индекс
	_, err = ind.Build(testInputData)
	if err != nil {
		t.Errorf("Ожидается: нет ошибки, получено: %v", err)
	}
	// Проверяем соответствие ключей в индексе
	expKeys := []string{"go", "googl", "ogle", "gl", "goo", "google", "oogle", "og", "ogl", "wh", "why", "hy", "goog", "oo", "le", "oog", "oogl", "gle"}
	for _, key := range expKeys {
		if ind.Hash[key] == nil {
			t.Errorf("Ключ %v не найден в индексе", key)
		}
	}
	err = ind.SaveData()
	if err != nil {
		t.Errorf("[build] %s", err)
	}
	// Проверяем соответвие ключей в индексе в записанном файле
	writedData, err := ReadData(&storage)
	if err != nil {
		t.Errorf("Не удается прочитать выходной файл")
	}
	for _, key := range expKeys {
		if writedData.Hash[key] == nil {
			t.Fatalf("Ключ %v не найден в индексе в выходном файле", key)
		}
	}
}
