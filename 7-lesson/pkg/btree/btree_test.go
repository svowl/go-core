package btree

import (
	"reflect"
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
	if root.Right.ID != a.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Value, a.Value, root.Right.ID, a.ID)
	}
	// В дереве root и a, root.Count ожидается 1 (root не учитывается)
	expCount := 1
	if root.Count != expCount {
		t.Fatalf("Count(1): получено %d, ожидается %d", root.Count, expCount)
	}
	root.Add(&b)
	if root.Left.ID != b.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Value, b.Value, root.Left.ID, b.ID)
	}
	root.Add(&c)
	if a.Right.ID != c.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Value, c.Value, a.Right.ID, c.ID)
	}
	root.Add(&d)
	if b.Left.ID != d.ID {
		t.Fatalf("В %s добавляем %s: получено %d, ожидается %d", root.Value, d.Value, b.Left.ID, d.ID)
	}
	// В дереве root, a, b, c, d - root.Count ожидается 4
	expCount = 4
	if root.Count != expCount {
		t.Fatalf("Count(2): получено %d, ожидается %d", root.Count, expCount)
	}
	e := Tree{
		ID:    5,
		Value: "e",
	}
	root.Add(&e)
	if d.Value != e.Value {
		t.Fatalf("В %s добавляем %s: получено %s, ожидается %s", root.Value, e.Value, d.Value, e.Value)
	}
	// e заменил d, поэтому root.Count не должно измениться
	if root.Count != expCount {
		t.Fatalf("Count(3): получено %d, ожидается %d", root.Count, expCount)
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

func TestTree_TreeMap(t *testing.T) {
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
	testTreeMap = []string{}
	exp := []string{"root", "b", "d", "a", "c"}
	root.TreeMap(callback)
	if !reflect.DeepEqual(testTreeMap, exp) {
		t.Fatalf("получено %v, ожидается %v", testTreeMap, exp)
	}
}

var testTreeMap []string

func callback(t *Tree) {
	testTreeMap = append(testTreeMap, t.Value.(string))
}
