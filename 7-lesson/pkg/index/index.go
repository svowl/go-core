package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go-core/7-lesson/pkg/btree"
	"go-core/7-lesson/pkg/crawler"
	"go-core/7-lesson/pkg/storage"
)

// Index это основная структура пакета
type Index struct {
	// Records это корневой узел бинарного дерева, содержащего список документов
	Records *btree.Tree
	// Обратный индекс
	Hash
	// Служебные поля
	Service
}

// Hash это структура данных для хранения обратного индекса
type Hash map[string][]int

// Service это сструктура, содержащая служебные данные
type Service struct {
	// Минимальная длина фрагмента (ключевого слова для поиска)
	fragmentMinLen int
	// Сервисный массив с URL в качестве ключей для проверки уникальности URL страниц
	urlHash hashURL
	// ids используем для контроля уникальности ID записей
	idHash hashID
	// storage принимает интерфейс storage.Interface для чтения и записи данных
	storage storage.Interface
}

// hashURL это структура для хранения хеша URL
type hashURL map[string]bool

// hashID это структура для хранения хеша ID
type hashID map[int]bool

// tmpData это структура для построения хешей URL и ID.
type tmpData struct {
	urlHash hashURL
	idHash  hashID
}

// tmp используется для генерации хешей при чтении данных из файла.
// Т.к. при обходе дерева вызывается callback-функция, мы используем временную глобальную переменную
var tmp tmpData

// FileData это структура, записываемая в файл.
// В файл скидываем только индекс и записи.
type FileData struct {
	Hash
	Records *btree.Tree
}

// New создает и инициализирует объект Index
func New() *Index {
	var i Index
	return &i
}

// Init создает и инициализирует объект Index
func (i *Index) Init(storage storage.Interface) error {
	i.Hash = make(map[string][]int)
	i.fragmentMinLen = 2
	i.urlHash = make(hashURL)
	i.idHash = make(hashID)
	i.storage = storage
	err := i.initData()
	if err != nil {
		return err
	}
	return nil
}

// Build получает на вход данные от spider'а, конвертирует их в Records и строит обратный индекс
func (i *Index) Build(data []crawler.Document) (Hash, error) {
	dataLen := len(data)
	ri := 0 // это ненужный в реальной системе артефакт, здесь используется для генерации случайного числа
	for _, doc := range data {
		// Проверяем - возможно этот URL уже обработан
		if _, found := i.urlHash[doc.URL]; found {
			continue
		}
		// id генерим случайным образом
		id, err := i.rnd(ri, dataLen)
		if err != nil {
			return nil, err
		}
		doc.ID = id
		node := btree.Tree{
			ID:    id,
			Count: 1,
			Value: doc,
		}
		if i.Records == nil {
			// Создаем корневой элемент дерева
			i.Records = &node
		} else {
			// Добавляем запись в список
			i.Records.Add(&node)
		}
		// Строим по ней обратный хеш с ключами - фрагментами слов из Record.Title
		for _, key := range fragments(doc.Title, i.fragmentMinLen) {
			i.Hash[key] = append(i.Hash[key], node.ID)
		}
		// Сохраняем URL и ID в хэш
		i.urlHash[doc.URL] = true
		i.idHash[id] = true
		ri++
	}
	return i.Hash, nil
}

// Search выдает список ID документов по строке
func (i *Index) Search(phrase string) ([]int, bool) {
	result, found := i.Hash[strings.ToLower(phrase)]
	if found != true {
		return nil, false
	}
	return result, true
}

// initData инициализирует начальные данные из файла
func (i *Index) initData() error {
	fileData, err := ReadData(i.storage)
	if err != nil {
		return err
	}
	if fileData.Hash == nil {
		return nil
	}
	i.Hash = fileData.Hash
	i.Records = fileData.Records
	// сбрасываем tmp в исходное состояние
	tmp = tmpData{
		make(hashURL),
		make(hashID),
	}
	// Инициализируем хеши URL и ID, они нужны для построения индекса
	i.Records.TreeMap(initHash)
	i.urlHash = tmp.urlHash
	i.idHash = tmp.idHash
	return nil
}

// initHash инициализирует хеши, т.к. URL и ID должны быть уникальны
func initHash(t *btree.Tree) {
	tmp.addItem(t)
}

// addItem выполняет добавление данных во временные хеши
func (tmp *tmpData) addItem(t *btree.Tree) {
	tmp.urlHash[t.Value.(crawler.Document).URL] = true
	tmp.idHash[t.ID] = true
}

// ReadData записывает данные индекса в файл
func ReadData(r storage.Interface) (FileData, error) {
	var fileData FileData
	data, err := r.Get()
	if err != nil {
		return fileData, err
	}
	if len(data) == 0 {
		// Файл пустой - допустимая ситуация, выходим без ошибки
		return fileData, nil
	}
	// Десериализуем data
	err = json.Unmarshal(data, &fileData)
	if err != nil {
		return fileData, err
	}
	// Корректируем Value после десериализации
	fileData.Records.TreeMap(convertTreeValue)
	return fileData, nil
}

// SaveData записывает данные индекса в файл
func (i *Index) SaveData() error {
	data := FileData{
		i.Hash,
		i.Records, //.Serialize(),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = i.storage.Set(jsonData)
	if err != nil {
		return err
	}
	return nil
}

// convertTreeValue приводит тип map[string]interface{} к типу index.Record,
// поскольку json.Unmarshal десериализует Records в виде структуры вида map[string]interface{}
func convertTreeValue(t *btree.Tree) {
	var value crawler.Document
	for k, v := range t.Value.(map[string]interface{}) {
		switch k {
		case "URL":
			value.URL = fmt.Sprintf("%s", v)
		case "Title":
			value.Title = fmt.Sprintf("%s", v)
		case "Body":
			value.Body = fmt.Sprintf("%s", v)
		}
	}
	t.Value = value
}

// fragments разбивает строку text на слова + все возможные фрагменты слов (Google: Go, Goo, gle...) длиной не менее fragmentMinLen символов.
// Все фрагменты в нижнем регистре. Комбинации слов не рассматриваем.
func fragments(text string, fragmentMinLen int) []string {
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
func (i *Index) rnd(ri int, max int) (int, error) {
	var id int
	lim := 100
	for {
		lim--
		rand.Seed(time.Now().UnixNano() * int64(ri+lim))
		id = rand.Intn(max * 1000 * lim)
		if _, isUsed := i.idHash[id]; !isUsed || lim <= 0 {
			break
		}
	}
	if lim >= 100 || id == 0 {
		return 0, errors.New("Ошибка генерации уникального индекса документа")
	}
	return id, nil
}
