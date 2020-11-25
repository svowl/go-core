package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var score []int

func player(ch chan string, n int, wg *sync.WaitGroup) {
	var newMsg string

	for {
		msg := <-ch

		// если пришло сообщение stop, ничего не делаем, продолжаем чтение канала дальше
		if msg == "stop" {
			continue
		}

		// С вероятностью 20% (1/5) игрок выигрывает подачу
		// Увеличиваем счет, выводим сообщение и уменьшаем счетчик wg
		if rand.Intn(5) == 4 {
			score[n]++
			fmt.Printf("Player #%d wins the serve: %v\n", n+1, score)
			ch <- "stop"
			wg.Done()

		} else {
			// Иначе продолжаем розыгрыш
			if msg == "ping" {
				// На ping отвечаем pong
				newMsg = "pong"
			} else {
				// На pong отвечаем ping
				newMsg = "ping"
			}
			fmt.Printf("Player #%d: %s\n", n+1, newMsg)
			ch <- newMsg
		}

		// Задержка, чтобы не мельтешило
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ch := make(chan string)
	var wg sync.WaitGroup
	score = []int{0, 0}

	go player(ch, 0, &wg)
	go player(ch, 1, &wg)

	for {
		wg.Add(1)
		fmt.Println("begin")
		// запуск подачи
		ch <- "begin"
		// ожидание выигрыша подачи одним из игроков
		wg.Wait()
		// Проверяем счет и прекращаем игру, если один из игроков набрал 3 очка
		if score[0] == 3 || score[1] == 3 {
			break
		}
		// Задержка между подачами полсекунды
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println("Final score:", score)
}
