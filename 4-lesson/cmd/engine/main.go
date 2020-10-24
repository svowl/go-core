package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"go-core/4-lesson/pkg/index"
	"go-core/4-lesson/pkg/spider"
)

// Scanner interface
type Scanner interface {
	Scan(string, int) (map[string]string, error)
}

// Scan method of Spider type
func Scan(s Scanner, url string, depth int) (map[string]string, error) {
	return s.Scan(url, depth)
}

func main() {

	// Берем поисковую фразу из командной строки
	var phrase string = ""
	flag.StringVar(&phrase, "s", "", "Укажите поисковую фразу")
	flag.Parse()

	fmt.Println("Сканируем...")
	s := new(spider.Spider)

	urls := []string{"https://go.dev", "https://www.google.com", "https://habr.com"}

	for _, url := range urls {
		fmt.Printf("  %s...\n", url)
		data, err := Scan(s, url, 2)
		if err != nil {
			log.Fatalf("ошибка при сканировании сайта %s: %v\n", url, err)
		}
		fmt.Printf("  ...найдено %d документов, индексируем...\n", len(data))
		// Строим индекс по списку просканированных документов
		index.Build(data)
	}

	// Сортируем список просканированных документов
	index.SortRecords()

	fmt.Printf("Проиндексировано %d страниц\n", len(index.Records))

	// Реализуем ввод фразы с клавиатуры и поиск в индексе
	reader := bufio.NewReader(os.Stdin)
	for {
		if phrase != "" {
			fmt.Printf("Поиск по строке \"%s\"\n", phrase)
			found := false
			for u, title := range search(phrase) {
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

// search находит записи в индексе по фразе.
func search(phrase string) map[string]string {
	return index.Search(phrase)
}
