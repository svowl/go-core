package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"go-core/3-lesson/pkg/spider"
)

// Scanner интерфейс, объявляющий метод Scan, реализуемый в пакетах spider и mem/spider
type Scanner interface {
	Scan(string, int) (map[string]string, error)
}

// scan запускает сканер s для сайтов urls и отдает результат сканирования
func scan(s Scanner, urls []string) map[string]string {
	storage := make(map[string]string)
	for _, url := range urls {
		fmt.Printf("  %s...\n", url)
		data, err := s.Scan(url, 2)
		if err != nil {
			log.Printf("ошибка при сканировании сайта %s: %v\n", url, err)
			continue
		}
		for k, v := range data {
			storage[k] = v
		}
	}
	return storage
}

// search находит записи в storage по фразе phrase,
// вынесено в отдельную функцию, чтобы протестировать непосредственно поиск
func search(phrase string, storage map[string]string) map[string]string {
	res := make(map[string]string)
	for u, title := range storage {
		if strings.Contains(strings.ToLower(title), strings.ToLower(phrase)) {
			res[u] = title
		}
	}
	return res
}

func main() {

	// Берем поисковую фразу из командной строки
	var phrase string = ""
	flag.StringVar(&phrase, "s", "", "Укажите поисковую фразу")
	flag.Parse()

	fmt.Println("Индексируем...")

	s := new(spider.Spider)

	// Структура для хранения данных
	// key: url, value: page title
	var storage map[string]string

	var urls = []string{"https://go.dev", "https://www.google.com"} //, "https://habr.com"}

	// Получаем данные сканирования
	storage = scan(s, urls)

	fmt.Printf("Теперь в индексе %d записей\n", len(storage))

	// Реализуем ввод фразы с клавиатуры
	for {
		if phrase != "" {
			fmt.Printf("Поиск по строке \"%s\"\n", phrase)
			found := false
			for u, title := range search(phrase, storage) {
				fmt.Printf("  %s: %s\n", u, title)
				found = true
			}
			if !found {
				fmt.Println("Ничего не найдено")
			}
		}
		fmt.Print("Введите поисковую фразу: ")
		fmt.Scanln(&phrase)
	}
}
