package builder

import (
	"errors"
	"testing"

	"go-core/7-lesson/pkg/index"
	memspider "go-core/7-lesson/pkg/spider/mem"
	memstorage "go-core/7-lesson/pkg/storage/mem"
)

func TestService_Build(t *testing.T) {
	b := New()
	sp := memspider.New()
	// Инициализируем с ошибкой
	stErr := memstorage.New()
	stErr.Error = errors.New("Test error")
	err := b.Init([]string{""}, sp, stErr)
	if err == nil {
		t.Fatalf("Ожидаем ошибку")
	}
	// Инициализируем без ошибки
	st := memstorage.New()
	err = b.Init([]string{""}, sp, st)
	if err != nil {
		t.Fatalf("Ошибка инициализации: %v", err)
	}
	// Запускаем билдер
	b.Build()
	// Проверяем соответвие ключей в индексе в записанном файле
	writedData, err := index.ReadData(st)
	if err != nil {
		t.Errorf("Не удается прочитать файл с результатом работы builder'а")
	}
	expKeys := []string{"go", "googl", "ogle", "gl", "goo", "google", "oogle", "og", "ogl", "wh", "why", "hy", "goog", "oo", "le", "oog", "oogl", "gle"}
	for _, key := range expKeys {
		if writedData.Hash[key] == nil {
			t.Fatalf("Ключ %v не найден в индексе в выходном файле", key)
		}
	}
}
