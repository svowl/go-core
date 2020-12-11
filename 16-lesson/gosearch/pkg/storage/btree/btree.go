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

// DeleteDoc удаляет документ из структуры, возвращает true при успехе
func (tree *Tree) DeleteDoc(id int) bool {
	if tree.Root == nil {
		// пустое дерево
		return false
	}
	if tree.Root.Doc.ID == id && tree.Root.Left == nil && tree.Root.Right == nil {
		// в дереве один единственный корневой узел и у него совпал ID
		tree.Root = nil
		return true
	}
	if tree.Root.Doc.ID == id {
		// ID совпал с корневой вершиной
		tmpRoot := &Node{Left: tree.Root}
		if delete(tree.Root, tmpRoot, id) {
			tree.Root = tmpRoot.Left
			return true
		}
		return false
	}
	return delete(tree.Root, nil, id)
}

// Count возвращает кол-во узлоа в дереве.
// Для оптимизации результат можно кэшировать, но это здесь не реализовано.
func (tree *Tree) Count() int {
	if tree.Root == nil {
		return 0
	}
	len := 0
	tree.Root.traverse(&len)
	return len
}

// All возвращает документы всех узлов дерева
func (tree *Tree) All() []crawler.Document {
	if tree.Root == nil {
		return nil
	}
	docs := make([]crawler.Document, 0, 10)
	tree.Root.all(&docs)
	return docs
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

// Delete реализует удаление узла из дерева по ID
func delete(node *Node, parent *Node, ID int) bool {
	if ID < node.Doc.ID {
		if node.Left == nil {
			return false
		}
		// Удаляем в левой части
		return delete(node.Left, node, ID)
	}
	if ID > node.Doc.ID {
		if node.Right == nil {
			return false
		}
		// Удаляем в правой части
		return delete(node.Right, node, ID)
	}
	// ID совпал с node.Doc.ID
	if node.Left != nil && node.Right != nil {
		// Есть оба дочерних узла, ищем минимальный в правом дереве и замещаем им текущий узел
		successor := node.Right.findMin()
		node.Doc = successor.Doc
		// удаляем successor из правого поддерева узла node
		return delete(node.Right, node, successor.Doc.ID)
	}
	if node.Left != nil {
		// Есть только левый дочерний узел - заменяем node на него
		replaceParent(node, parent, node.Left)
		return true
	}
	if node.Right != nil {
		// Есть только правый дочерний узел - заменяем node на него
		replaceParent(node, parent, node.Right)
		return true
	}
	// Нет дочерних узлов у node - удаляем в parent ссылку на него
	replaceParent(node, parent, nil)
	return true
}

// replaceParent изменяет Left/Right у parent так, чтобы они указывали на новый узел
func replaceParent(node *Node, parent *Node, newNode *Node) {
	if parent != nil {
		if node == parent.Left {
			parent.Left = newNode
		} else {
			parent.Right = newNode
		}
	}
	if newNode != nil {
		node.Right = newNode
	}
}

// findMin возвращает узел с минимальным ID поддереве начиная с node
func (node *Node) findMin() *Node {
	if node.Left == nil {
		return node
	}
	return node.Left.findMin()
}

// traverse обходит каждый узел дерева, в каждом узле прибавляет 1 в переданный счетчик.
// Используется для вычисления количества документов в дереве, см. tree.Count()
func (node *Node) traverse(c *int) {
	if node == nil {
		return
	}
	*c++
	node.Left.traverse(c)
	node.Right.traverse(c)
}

// all обходит каждый узел дерева и записывает документ в массив docs
func (node *Node) all(docs *[]crawler.Document) {
	if node == nil {
		return
	}

	node.Left.all(docs)
	*docs = append(*docs, node.Doc)
	node.Right.all(docs)
}
