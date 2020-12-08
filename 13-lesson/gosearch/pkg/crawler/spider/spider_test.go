package spider

import (
	"gosearch/pkg/crawler"
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestService_BatchScan(t *testing.T) {

	if !testing.Short() {
		t.Skip("только с флагом `-short`")
	}

	// Тестовая страница заведена специально для тестирования пакета spider,
	// имеет постоянную структуру и никогда не меняется.
	urls := []string{"http://svowl.github.io/test/"}

	s := New()

	var got []crawler.Document
	var wg sync.WaitGroup
	wg.Add(1)

	chRes, chErr := s.BatchScan(urls, 2, 1)
	go func() {
		defer wg.Done()
		for doc := range chRes {
			got = append(got, doc)
		}
	}()
	go func() {
		for err := range chErr {
			t.Fatal("Ошибка сканирования: ", err)
		}
	}()

	wg.Wait()

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

func TestService_Scan(t *testing.T) {

	if !testing.Short() {
		t.Skip("только с флагом `-short`")
	}

	// Тестовая страница заведена специально для тестирования пакета spider,
	// имеет постоянную структуру и никогда не меняется.
	url := "http://svowl.github.io/test/"

	s := New()

	got, err := s.Scan(url, 2)
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
