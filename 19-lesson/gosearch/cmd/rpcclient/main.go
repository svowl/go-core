package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"

	"gosearch/pkg/crawler"
)

// Устанавливает соединение, запрашивает результаты поиска по phrase и возвращает их в текстовом виде
func search(phrase string) (string, error) {
	// Создаем соединение
	client, err := rpc.DialHTTPPath("tcp4", "localhost:80", "/rpc")
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Выполняем запрос документов по фразе
	var req = &phrase
	var data []crawler.Document
	client.Call("Service.Search", req, &data)
	if err != nil {
		return "", err
	}

	// Обрабатываем результат и выдаем его в виде строки
	msg := ""
	for _, doc := range data {
		msg = msg + fmt.Sprintf("%s: %s\r\n", doc.URL, doc.Title)
	}

	if msg == "" {
		msg = "Нет результатов"
	}

	return msg, nil
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nВведите поисковую фразу: ")
		phrase, _ := reader.ReadString('\n')
		phrase = strings.Replace(phrase, "\r\n", "", -1)
		phrase = strings.Replace(phrase, "\n", "", -1)

		if phrase != "" {
			response, err := search(phrase)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Результат:")
			fmt.Println(response)
		}
	}
}
