package index

import (
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// Record содержит запись об одном просканированном документе
type Record struct {
	ID    int
	URL   string
	Title string
}

// Records содержит список просканированных документов
var Records []Record

// ByID реализует sort.Interface для []Record, сортировка по Record.ID
type ByID []Record

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Hash хранит обратный индекс
var Hash map[string][]int = make(map[string][]int)

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
		record := Record{
			ID:    id,
			URL:   u,
			Title: t,
		}
		// Добавляем запись в список
		Records = append(Records, record)
		// Строим по ней обратный хеш с ключами - фрагментами слов из Record.Title
		for _, key := range fragments(record.Title) {
			Hash[key] = append(Hash[key], record.ID)
		}
		// Сохраняем URL и ID в хэш
		urlHash[u] = true
		idHash[id] = true
		i++
	}
	return Hash
}

// Search находит проиндексированные записи по фразе,
func Search(phrase string) map[string]string {
	res := make(map[string]string)
	if ids, found := Hash[strings.ToLower(phrase)]; found {
		// Фраза найдена в хеше, ids содержит индексы документов (Record.ID) в массиве Records
		var total = len(Records) - 1
		for _, id := range ids {
			// Используем бинарный поиск из стандартного пакета, recID - номер записи в Records
			recID := sort.Search(total, func(i int) bool { return Records[i].ID >= id })
			res[Records[recID].URL] = Records[recID].Title
		}
	}
	return res
}

// sortRecords сортирует список документов Records
func SortRecords() []Record {
	sort.Sort(ByID(Records))
	return Records
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
