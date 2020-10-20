package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"go-core/3-lesson/pkg/spider"
)

// Scanner interface
type Scanner interface {
	Scan(string, int) (map[string]string, error)
}

// Spider type
type Spider struct{}

// Scan method of Spider type
func (*Spider) Scan(url string, depth int) (map[string]string, error) {
	return spider.Scan(url, depth)
}

func main() {

	// Берем поисковую фразу из командной строки
	var phrase string = ""
	flag.StringVar(&phrase, "s", "", "Укажите поисковую фразу")
	flag.Parse()

	fmt.Println("Индексируем...")

	s := new(Spider)

	// Структура для хранения данных
	// key: url, value: page title
	storage := make(map[string]string)

	urls := []string{"https://go.dev", "https://www.google.com"} //, "https://habr.com"}
	for _, url := range urls {
		fmt.Printf("  %s...\n", url)
		data, err := s.Scan(url, 2)
		if err != nil {
			log.Printf("ошибка при сканировании сайта %s: %v\n", url, err)
		}

		for k, v := range data {
			storage[k] = v
		}
	}

	fmt.Printf("Теперь в индексе %d записей\n", len(storage))

	// Реализуем ввод фразы с клавиатуры
	reader := bufio.NewReader(os.Stdin)
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
		phrase, _ = reader.ReadString('\n')
		phrase = strings.Replace(phrase, "\r\n", "", -1)
		phrase = strings.Replace(phrase, "\n", "", -1)
	}
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
