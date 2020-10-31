package btree

// Tree описывает узел дерева
type Tree struct {
	ID          int
	left, right *Tree
	Value       interface{}
}

// Add добавляет узел в дерево
func (b *Tree) Add(n *Tree) {
	if b.ID == n.ID {
		// Ключ совпал, переписываем значение
		b.Value = n.Value
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
func (b *Tree) Search(ID int) *Tree {
	if ID == b.ID {
		// ID совпадает, возвращаем элемент
		return b
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
	return nil
}
