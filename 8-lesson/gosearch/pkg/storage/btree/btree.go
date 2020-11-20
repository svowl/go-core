package btree

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"gosearch/pkg/crawler"
)

// Tree структура, хранящая бинарное дерево элементов
type Tree struct {
	Root *Node
}

// Node описывает узел дерева
type Node struct {
	Left, Right *Node
	Doc         crawler.Document
}

// New создает объект Tree
func New() *Tree {
	var tree Tree
	// Инициализация генератора случайных чисел для генерации ID документов
	rand.Seed(time.Now().UnixNano())
	return &tree
}

// AddDoc добавляет документ в дерево
func (tree *Tree) AddDoc(doc crawler.Document) {
	node := Node{
		Doc: doc,
	}
	if tree.Root == nil {
		// Добавляем в корень
		tree.Root = &node
	} else {
		// Добавляем в дерево
		tree.Root.Add(&node)
	}
}

// FindDoc ищет документ в структуре. Возвращает документ и флаг - статус поиска.
func (tree *Tree) FindDoc(id int) (crawler.Document, bool) {
	var doc crawler.Document
	if tree.Root == nil {
		return doc, false
	}
	node := tree.Root.Search(id)
	if node == nil {
		var doc crawler.Document
		return doc, false
	}
	return node.Doc, true
}

// Count возвращает кол-во узлоа в дереве.
// Для оптимизации результат можно кэшировать, но это здесь не реализовано.
func (tree *Tree) Count() int {
	if tree.Root == nil {
		return 0
	}
	ch := make(chan int, 1)
	len := 1
	go func(ch <-chan int, len *int) {
		for {
			select {
			case <-ch:
				*len++
			}
		}
	}(ch, &len)
	tree.Root.traverse(ch)
	close(ch)
	return len
}

// GenerateID генерирует уникальный ID документа.
// Для бинарного дерева принципиально важно, чтобы ID документов не были последовательными,
// поэтому здесь используется генерация случайного числа.
// Понятно, что это грубая имитация, чтобы не тормозить работу программы.
// Todo: В качестве оптимизации на будущее: генерировать числовой ID из значения doc.URL, чтобы предотвратить дубликаты.
func (tree *Tree) GenerateID() (int, error) {
	var ID, tries int
	max := (tree.Count() + 1) * 10_000
	for ID == 0 && tries < 10 {
		ID = rand.Intn(max)
		tries++
		if _, found := tree.FindDoc(ID); found {
			continue
		}
		break
	}
	if ID == 0 {
		return 0, errors.New("Ошибка генерации уникального ID")
	}
	return ID, nil
}

// JSONData возвращает структуру в сериализованном виде
func (tree *Tree) JSONData() ([]byte, error) {
	jsonData, err := json.Marshal(tree)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// LoadFromJSON загружает в дерево данные из json-строки
func (tree *Tree) LoadFromJSON(data []byte) error {
	err := json.Unmarshal(data, tree)
	if err != nil {
		return err
	}
	return nil
}

// Add добавляет узел в дерево, возвращает 1, если узел добавлен, 0 - если переписан существующий (при совпадении ID).
// Возвращаемый результат сейчас не используется, осталось от прошлой реализации. Может пригодиться при использовании кэша Count(), поэтому оставил.
func (node *Node) Add(n *Node) int {
	if node.Doc.ID == n.Doc.ID {
		// Ключ совпал, переписываем значение
		node.Doc = n.Doc
		return 0
	}
	if n.Doc.ID < node.Doc.ID {
		// Нужно добавить элемент в левую часть
		if node.Left != nil {
			// Левый элемент уже есть, выполняем добавление к нему
			return node.Left.Add(n)
		}
		// Записываем ссылку на новый элемент в Left
		node.Left = n
		return 1
	}
	// Добавляем в правую ветку
	if node.Right != nil {
		// Правый элемент уже есть, добавляем к нему
		return node.Right.Add(n)
	}
	// Записываем ссылку на новый элемент в Right
	node.Right = n
	return 1
}

// Search реализует рекурсивный поиск узла в дереве по ID
func (node *Node) Search(ID int) *Node {
	if ID == node.Doc.ID {
		// ID совпадает, возвращаем элемент
		return node
	}
	if ID < node.Doc.ID && node.Left != nil {
		// Ищем в левой части
		return node.Left.Search(ID)
	}
	if ID > node.Doc.ID && node.Right != nil {
		// Ищем в правой части
		return node.Right.Search(ID)
	}
	// Не нашли - возвращаем nil
	return nil
}

// traverse обходит каждый узел дерева, в каждом узле записывает 1 в переданный канал
// Используется для вычисления количества документов в дереве, см. tree.Count()
func (node *Node) traverse(ch chan<- int) {
	if node == nil {
		return
	}
	ch <- 1
	node.Left.traverse(ch)
	node.Right.traverse(ch)
}
