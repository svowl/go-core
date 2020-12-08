package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// Устанавливает соединение, запрашивает результаты поиска по phrase и возвращает их в текстовом виде
func connect(phrase string) (string, error) {
	// Создаем соединение
	conn, err := net.Dial("tcp4", "localhost:8000")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Таймаут на чтение полсекунды
	conn.SetDeadline(time.Now().Add(time.Millisecond * 500))

	// Пишем в соединение поисковую фразу
	_, err = conn.Write([]byte(phrase + "\n"))
	if err != nil {
		return "", err
	}

	// Читаем ответ, записываем его в msg
	msg := make([]byte, 0, 4096)
	tmp := make([]byte, 20)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			break
		}
		msg = append(msg, tmp[:n]...)
	}

	return string(msg), nil
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nВведите поисковую фразу: ")
		phrase, _ := reader.ReadString('\n')
		phrase = strings.Replace(phrase, "\r\n", "", -1)
		phrase = strings.Replace(phrase, "\n", "", -1)

		if phrase != "" {
			response, err := connect(phrase)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Результат:")
			fmt.Println(response)
		}
	}
}
