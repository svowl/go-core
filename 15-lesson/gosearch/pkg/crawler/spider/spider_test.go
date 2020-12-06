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
	urls := []string{"http://svowl.github.io/test/"}

	s := New(1)

	got, err := s.Scan(urls, 2)
	if err != nil {
		t.Fatal(err)
	}

	// Сортировка нужна для точного сравнения с помощью reflect.DeepEqual(), т.к. Scan возвращает несортированные данные
	sort.Slice(got, func(i, j int) bool { return got[i].Title < got[j].Title })

	want := []crawler.Document{
		{URL: "http://svowl.github.io/test/", Title: "Test 1"},
		{URL: "http://svowl.github.io/test2.html", Title: "Test 2"},
		{URL: "http://svowl.github.io/test3.html", Title: "Test 3"},
		{URL: "http://svowl.github.io/test4.html", Title: "Test 4"},
		{URL: "http://svowl.github.io/test/test/test5.html", Title: "Test 5"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Получено: %v, ожидается: %v", got, want)
	}
}
