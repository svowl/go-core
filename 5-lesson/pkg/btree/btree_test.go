package btree

import (
	"testing"
)

func TestNode_Add(t *testing.T) {
	root := Node{
		ID:  10,
		Doc: "root",
	}
	a := Node{
		ID:  12,
		Doc: "a",
	}
	b := Node{
		ID:  8,
		Doc: "b",
	}
	c := Node{
		ID:  18,
		Doc: "c",
	}
	d := Node{
		ID:  5,
		Doc: "d",
	}
	root.Add(&a)
	if root.right.ID != a.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Doc, a.Doc, root.right.ID, a.ID)
	}
	root.Add(&b)
	if root.left.ID != b.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Doc, b.Doc, root.left.ID, b.ID)
	}
	root.Add(&c)
	if a.right.ID != c.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Doc, c.Doc, a.right.ID, c.ID)
	}
	root.Add(&d)
	if b.left.ID != d.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Doc, d.Doc, b.left.ID, d.ID)
	}
	e := Node{
		ID:  5,
		Doc: "e",
	}
	root.Add(&e)
	if d.Doc != e.Doc {
		t.Fatalf("В %s добавляем %s: получено %s, ожидается %s", root.Doc, e.Doc, d.Doc, e.Doc)
	}

}

func TestNode_Search(t *testing.T) {
	root := Node{
		ID:  10,
		Doc: "root",
	}
	a := Node{
		ID:  12,
		Doc: "a",
	}
	b := Node{
		ID:  8,
		Doc: "b",
	}
	c := Node{
		ID:  18,
		Doc: "c",
	}
	d := Node{
		ID:  5,
		Doc: "d",
	}
	root.Add(&a)
	root.Add(&b)
	root.Add(&c)
	root.Add(&d)
	// Проверка поиска существующих значений
	ids := []int{18, 8, 5, 10, 12}
	for _, ID := range ids {
		got, err := root.Search(ID)
		if err != nil {
			t.Fatalf("Не найден элемент %d", ID)
		}
		if got.ID != ID {
			t.Fatalf("Получено %d, ожидается %d", got.ID, ID)
		}
	}
	// Проверка поиска на несуществующее значение
	_, err := root.Search(200)
	if err == nil {
		t.Fatalf("Найдено несуществующее значение %d", 200)
	}
}
