package btree

import (
	"testing"
)

func TestTree_Add(t *testing.T) {
	root := Tree{
		ID:    10,
		Value: "root",
	}
	a := Tree{
		ID:    12,
		Value: "a",
	}
	b := Tree{
		ID:    8,
		Value: "b",
	}
	c := Tree{
		ID:    18,
		Value: "c",
	}
	d := Tree{
		ID:    5,
		Value: "d",
	}
	root.Add(&a)
	if root.right.ID != a.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Value, a.Value, root.right.ID, a.ID)
	}
	root.Add(&b)
	if root.left.ID != b.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Value, b.Value, root.left.ID, b.ID)
	}
	root.Add(&c)
	if a.right.ID != c.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Value, c.Value, a.right.ID, c.ID)
	}
	root.Add(&d)
	if b.left.ID != d.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Value, d.Value, b.left.ID, d.ID)
	}
	e := Tree{
		ID:    5,
		Value: "e",
	}
	root.Add(&e)
	if d.Value != e.Value {
		t.Fatalf("В %s добавляем %s: получено %s, ожидается %s", root.Value, e.Value, d.Value, e.Value)
	}

}

func TestTree_Search(t *testing.T) {
	root := Tree{
		ID:    10,
		Value: "root",
	}
	a := Tree{
		ID:    12,
		Value: "a",
	}
	b := Tree{
		ID:    8,
		Value: "b",
	}
	c := Tree{
		ID:    18,
		Value: "c",
	}
	d := Tree{
		ID:    5,
		Value: "d",
	}
	root.Add(&a)
	root.Add(&b)
	root.Add(&c)
	root.Add(&d)
	// Проверка поиска существующих значений
	ids := []int{18, 8, 5, 10, 12}
	for _, ID := range ids {
		got := root.Search(ID)
		if got == nil {
			t.Fatalf("Не найден элемент %d", ID)
		}
		if got.ID != ID {
			t.Fatalf("Получено %d, ожидается %d", got.ID, ID)
		}
	}
	// Проверка поиска на несуществующее значение
	got := root.Search(200)
	if got != nil {
		t.Fatalf("Найдено несуществующее значение %d", 200)
	}
}
