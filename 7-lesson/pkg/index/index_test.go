package index

import (
	"errors"
	"go-core/7-lesson/pkg/storage/mem"
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
var testInputData = map[string]string{
	"https://google.com": "Google",
	"https://go.dev/":    "Why Go",
}

func TestNew(t *testing.T) {
	// Empty storage
	storage := mem.Storage{}
	_, err := New(&storage)
	if err != nil {
		t.Error(err)
	}
	// Storage с некорректными json-данными
	storage = mem.Storage{
		Content: []byte(testContentWrong),
	}
	_, err = New(&storage)
	if err == nil {
		t.Errorf("Ожидается ошибка json, получено nil")
	}
	// Storage с ошибкой
	testError := errors.New("Test error")
	storage = mem.Storage{
		Error: testError,
	}
	_, err = New(&storage)
	if err != testError {
		t.Errorf("Ожидается ошибка %v, получено %v", testError, err)
	}
	// Storage с корректной json-строкой
	storage = mem.Storage{
		Content: []byte(testContentCorrect),
	}
	_, err = New(&storage)
	if err != nil {
		t.Errorf("Ожидается: нет ошибки, получено: %v", err)
	}
}

func TestIndex_Build(t *testing.T) {
	// Создаем индекс с пустым исходным файлом
	storage := mem.Storage{}
	index, err := New(&storage)
	if err != nil {
		t.Error(err)
	}
	// Строим индекс
	_, err = index.Build(testInputData)
	if err != nil {
		t.Errorf("Ожидается: нет ошибки, получено: %v", err)
	}
	// Проверяем соответствие ключей в индексе
	var hashKeys []string
	for i, j := range index.Hash {
		hashKeys = append(hashKeys, i)
		_ = j
	}
	expKeys := []string{"go", "googl", "ogle", "gl", "goo", "google", "oogle", "og", "ogl", "wh", "why", "hy", "goog", "oo", "le", "oog", "oogl", "gle"}
	for _, key := range expKeys {
		if index.Hash[key] == nil {
			t.Errorf("Ключ %v не найден в индексе", key)
		}
	}
	err = index.SaveData()
	if err != nil {
		t.Errorf("[build] %s", err)
	}
	// Проверяем соответвие ключей в индексе в записанном файле
	writedData, err := ReadData(&storage)
	if err != nil {
		t.Errorf("Не удается прочитать выходной файл")
	}
	var hashKeys2 []string
	for i, j := range writedData.Hash {
		hashKeys2 = append(hashKeys2, i)
		_ = j
	}
	for _, key := range expKeys {
		if writedData.Hash[key] == nil {
			t.Fatalf("Ключ %v не найден в индексе в выходном файле", key)
		}
	}
}
