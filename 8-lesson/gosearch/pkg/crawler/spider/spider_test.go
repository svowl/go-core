package spider

import (
	"gosearch/pkg/crawler"
	"reflect"
	"sort"
	"testing"
)

func TestService_Scan(t *testing.T) {

	if !testing.Short() {
		t.Skip("только с флагом `-short`")
	}

	// Тестовая страница заведена специально для тестирования пакета spider,
	// имеет постоянную структуру и никогда не меняется.
	const url = "http://svowl.github.io/test/index.html"

	s := New()

	data, err := s.Scan(url, 2)
	if err != nil {
		t.Fatal(err)
	}

	// Сортировка нужна для точного сравнения с помощью reflect.DeepEqual(), т.к. Scan возвращает несортированные данные
	sort.Slice(data, func(i, j int) bool { return data[i].Title < data[j].Title })

	want := []crawler.Document{
		{URL: "http://svowl.github.io/test/index.html", Title: "Test 1"},
		{URL: "http://svowl.github.io/test2.html", Title: "Test 2"},
		{URL: "http://svowl.github.io/test3.html", Title: "Test 3"},
		{URL: "http://svowl.github.io/test4.html", Title: "Test 4"},
		{URL: "http://svowl.github.io/test/test/test5.html", Title: "Test 5"},
	}

	if !reflect.DeepEqual(data, want) {
		t.Fatalf("Получено: %v, ожидается: %v", data, want)
	}
}
