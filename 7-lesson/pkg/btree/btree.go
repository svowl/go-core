package btree

// Tree описывает узел дерева
// ID это индекс узла дерева
// Count это общее кол-во узлов в дереве начиная с текущего узла
// Value хранит произвольные данные
type Tree struct {
	ID          int
	Count       int
	Left, Right *Tree
	Value       interface{}
}

// Add добавляет узел в дерево
func (b *Tree) Add(n *Tree) int {
	if b.ID == n.ID {
		// Ключ совпал, переписываем значение
		b.Value = n.Value
		return 0
	}
	i := 1
	if n.ID < b.ID {
		// Нужно добавить элемент в левую часть
		if b.Left != nil {
			// Левый элемент уже есть, выполняем добавление к нему
			i = b.Left.Add(n)
		} else {
			// Записываем ссылку на новый элемент в Left
			b.Left = n
		}
		b.Count += i
		return i
	}
	// Добавляем в правую ветку
	if b.Right != nil {
		// Правый элемент уже есть, добавляем к нему
		i = b.Right.Add(n)
	} else {
		// Записываем ссылку на новый элемент в Right
		b.Right = n
	}
	b.Count += i
	return i
}

// Search реализует рекурсивный поиск узла в дереве по ID
func (b *Tree) Search(ID int) *Tree {
	if ID == b.ID {
		// ID совпадает, возвращаем элемент
		return b
	}
	if ID < b.ID && b.Left != nil {
		// Ищем в левой части
		return b.Left.Search(ID)
	}
	if ID > b.ID && b.Right != nil {
		// Ищем в правой части
		return b.Right.Search(ID)
	}
	// Не нашли - возвращаем nil
	return nil
}

// TreeMap обходит дерево и вызывает функцию callback для каждого узла
func (b *Tree) TreeMap(callback func(*Tree)) {
	if b == nil {
		return
	}
	callback(b)
	b.Left.TreeMap(callback)
	b.Right.TreeMap(callback)
}
