package index

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"go-core/5-lesson/pkg/btree"
)

// Hash хранит обратный индекс
var Hash map[string][]int = make(map[string][]int)

// Record это структура данных, хранящая информацию об одном документе
// ID вынесено в структуру узла btree.Node
type Record struct {
	URL   string
	Title string
}

// Records это корневой узел бинарного дерева, содержащего список документов
var Records *btree.Node

// RecordsCount счетчик количества документов (просто для вывода, чтобы не считать полным обходом дерева)
var RecordsCount int

// Минимальная длина фрагмента (ключевого слова для поиска)
const fragmentMinLen = 2

// Сервисный массив с URL в качестве ключей для проверки уникальности URL страниц
var urlHash = make(map[string]bool)

// ids используем для контроля уникальности ID записей
var idHash = make(map[int]bool)

// Build получает на вход данные от spider'а, конвертирует их в Records и строит обратный индекс
func Build(data map[string]string) map[string][]int {
	dataLen := len(data)
	i := 0 // это ненужный в реальной системе артефакт, здесь используется для генерации случайного числа
	for u, t := range data {
		// Проверяем - возможно этот URL уже обработан
		if _, found := urlHash[u]; found {
			continue
		}
		// id генерим случайным образом
		id := rnd(i, dataLen)
		node := btree.Node{
			ID: id,
			Doc: Record{
				URL:   u,
				Title: t,
			},
		}
		if Records == nil {
			// Создаем корневой элемент дерева
			Records = &node
		} else {
			// Добавляем запись в список
			Records.Add(&node)
		}
		RecordsCount++
		// Строим по ней обратный хеш с ключами - фрагментами слов из Record.Title
		for _, key := range fragments(node.Doc.(Record).Title) {
			Hash[key] = append(Hash[key], node.ID)
		}
		// Сохраняем URL и ID в хэш
		urlHash[u] = true
		idHash[id] = true
		i++
	}
	return Hash
}

// Search находит проиндексированные записи по фразе,
func Search(phrase string) []Record {
	var res []Record
	if ids, found := Hash[strings.ToLower(phrase)]; found {
		// Фраза найдена в хеше, ids содержит индексы документов (Record.ID) в массиве Records
		for _, id := range ids {
			// Используем бинарный поиск из стандартного пакета, recID - номер записи в Records
			record, _ := Records.Search(id)
			doc := record.Doc.(Record)
			res = append(res, doc)
		}
	}
	return res
}

// fragments разбивает строку text на слова + все возможные фрагменты слов (Google: Go, Goo, gle...) длиной не менее fragmentMinLen символов.
// Все фрагменты в нижнем регистре. Комбинации слов не рассматриваем.
func fragments(text string) []string {
	var words []string
	for _, word := range strings.Fields(strings.ToLower(text)) {
		if len(word) >= fragmentMinLen {
			for i := 0; i <= len(word)-fragmentMinLen; i++ {
				for j := i + fragmentMinLen; j <= len(word); j++ {
					words = append(words, word[i:j])
				}
			}
		}
	}
	return words
}

// Генерация уникального значения для Record.ID
func rnd(i int, max int) int {
	var id int
	lim := 100
	for {
		lim--
		rand.Seed(time.Now().UnixNano() * int64(i+lim))
		id = rand.Intn(max * 1000 * lim)
		if _, isUsed := idHash[id]; !isUsed || lim <= 0 {
			break
		}
	}
	if lim >= 100 || id == 0 {
		log.Fatalf("Ошибка генерации уникального индекса документа")
	}
	return id
}
