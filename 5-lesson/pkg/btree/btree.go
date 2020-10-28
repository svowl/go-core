package btree

import "errors"

// Node описывает узел дерева
type Node struct {
	ID          int
	left, right *Node
	Doc         interface{}
}

// Add добавляет узел в дерево
func (b *Node) Add(n *Node) {
	if b.ID == n.ID {
		// Ключ совпал, переписываем значение
		b.Doc = n.Doc
		return
	}
	if n.ID < b.ID {
		// Нужно добавить элемент в левую часть
		if b.left != nil {
			// Левый элемент уже есть, выполняем добавление к нему
			b.left.Add(n)
		} else {
			// Записываем ссылку на новый элемент в left
			b.left = n
		}
		return
	}
	// Добавляем в правую ветку
	if b.right != nil {
		// Правый элемент уже есть, добавляем к нему
		b.right.Add(n)
	} else {
		// Записываем ссылку на новый элемент в right
		b.right = n
	}
}

// Search реализует рекурсивный поиск узла в дереве по ID
func (b *Node) Search(ID int) (*Node, error) {
	if ID == b.ID {
		// ID совпадает, возвращаем элемент
		return b, nil
	}
	if ID < b.ID && b.left != nil {
		// Ищем в левой части
		return b.left.Search(ID)
	}
	if ID > b.ID && b.right != nil {
		// Ищем в правой части
		return b.right.Search(ID)
	}
	// Не нашли - возвращаем nil
	return nil, errors.New("Неуспешный поиск")
}
