package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"go-core/2-lesson/pkg/spider"
)

func main() {
	urls := []string{"https://go.dev", "https://www.google.com", "https://habr.com"}

	// Структура для хранения данных
	// key: url, value: page title
	storage := make(map[string]string)

	fmt.Println("Индексируем...")

	for _, url := range urls {
		fmt.Printf("  %s...\n", url)
		data, err := spider.Scan(url, 2)
		if err != nil {
			log.Printf("ошибка при сканировании сайта %s: %v\n", url, err)
		}

		for k, v := range data {
			storage[k] = v
		}
	}

	fmt.Printf("Теперь в индексе %d записей\n", len(storage))

	fmt.Println(storage)
	// Берем поисковую фразу из командной строки
	var phrase string = ""
	flag.StringVar(&phrase, "s", "", "Укажите поисковую фразу")
	flag.Parse()

	// Реализуем ввод фразы с клавиатуры
	reader := bufio.NewReader(os.Stdin)
	for {
		if phrase != "" {
			fmt.Printf("Поиск по строке \"%s\"\n", phrase)
			for u, title := range storage {
				if strings.Contains(strings.ToLower(title), strings.ToLower(phrase)) {
					fmt.Printf("  %s: %s\n", u, title)
				}
			}
		}
		fmt.Print("Введите поисковую фразу: ")
		phrase, _ = reader.ReadString('\n')
		phrase = strings.Replace(phrase, "\r\n", "", -1)
		phrase = strings.Replace(phrase, "\n", "", -1)
	}
}
