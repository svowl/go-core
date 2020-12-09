package btree

import (
	"gosearch/pkg/crawler"
	"reflect"
	"testing"
)

// Входные данные для тестирования
var docs = []crawler.Document{
	{ID: 10, URL: "https://google.com", Title: "Google"},
	{ID: 12, URL: "https://go.dev/", Title: "Why Go"},
	{ID: 8, URL: "https://golang.org/", Title: "The Go Programming Language"},
	{ID: 18, URL: "https://yandex.ru/", Title: "Yandex"},
	{ID: 5, URL: "https://rambler.ru/", Title: "Rambler"},
	{ID: 11, URL: "https://mail.ru/", Title: "Mail.ru"},
	{ID: 5, URL: "https://yahoo.com/", Title: "Yahoo!"},
}

func TestTree_AddDoc(t *testing.T) {
	tree := New()

	// Поиск в пустом дереве
	expID := 18
	doc, found := tree.FindDoc(expID)
	if found {
		t.Fatalf("Поиск %v: документ внезапно найден в пустом дереве", expID)
	}

	// В дереве пусто, проверяем Len
	expCount := 0
	gotCount := tree.Count()
	if gotCount != expCount {
		t.Fatalf("Count(0): получено %d, ожидается %d", gotCount, expCount)
	}

	// Добавляем первый узел (корневой), ID:10
	tree.AddDoc(docs[0])
	if tree.Root.Doc.ID != docs[0].ID {
		t.Fatalf("Добавление корневого узла %v: добавлен %v", docs[0].ID, tree.Root.Doc.ID)
	}
	if tree.Root.Right != nil {
		t.Fatalf("Добавление корневого узла: правый узел не nil")
	}
	if tree.Root.Left != nil {
		t.Fatalf("Добавление корневого узла: левый узел не nil")
	}

	// Добавляем узел 12:  10
	//                       \
	//                        12
	tree.AddDoc(docs[1])
	if tree.Root.Right.Doc.ID != docs[1].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", docs[1].ID, tree.Root.Right.Doc.ID)
	}

	// Добавляем узел 8:  10
	//                   /   \
	//                  8     12
	tree.AddDoc(docs[2])
	if tree.Root.Left.Doc.ID != docs[2].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", docs[2].ID, tree.Root.Left.Doc.ID)
	}

	// Добавляем узел 18:  10
	//                    /   \
	//                   8     12
	//                           \
	//                            18
	tree.AddDoc(docs[3])
	if tree.Root.Right.Right.Doc.ID != docs[3].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", docs[3].ID, tree.Root.Right.Right.Doc.ID)
	}

	// Добавляем узел 5:   10
	//                    /   \
	//                   8     12
	//                  /        \
	//                 5          18
	tree.AddDoc(docs[4])
	if tree.Root.Left.Left.Doc.ID != docs[4].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", docs[4].ID, tree.Root.Left.Left.Doc.ID)
	}

	// Добавляем узел 11:   10
	//                    /   \
	//                   8     12
	//                  /     /  \
	//                 5    11    18
	tree.AddDoc(docs[5])
	if tree.Root.Right.Left.Doc.ID != docs[5].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", docs[5].ID, tree.Root.Right.Left.Doc.ID)
	}

	// В дереве 6 элементов, проверяем Count
	expCount = 6
	gotCount = tree.Count()
	if gotCount != expCount {
		t.Fatalf("Count(1): получено %d, ожидается %d", gotCount, expCount)
	}

	// Добавляем документ с существующим ID (5), документ должен переписаться
	tree.AddDoc(docs[6])
	if tree.Root.Left.Left.Doc.Title != docs[6].Title {
		t.Fatalf("Добавление узла %v: добавлен %v", docs[6].Title, tree.Root.Left.Left.Doc.Title)
	}

	// В дереве должно остаться столько же элементов, проверяем Count
	gotCount = tree.Count()
	if gotCount != expCount {
		t.Fatalf("Count(2): получено %d, ожидается %d", gotCount, expCount)
	}

	// Тестируем получение списка всех документов
	gotAll := tree.All()
	if len(gotAll) != 6 {
		t.Errorf("получено %v, ожидается %v", len(gotAll), gotCount)
	}

	// Проверяем поиск в правой части
	expID = 18
	doc, found = tree.FindDoc(expID)
	if !found {
		t.Fatalf("Поиск %v: документ не найден", expID)
	}
	if doc.ID != expID {
		t.Fatalf("Поиск: получено %v, ожидается %v", doc.ID, expID)
	}

	// Проверяем поиск в левой части
	expID = 8
	doc, found = tree.FindDoc(expID)
	if !found {
		t.Fatalf("Поиск %v: документ не найден", expID)
	}
	if doc.ID != expID {
		t.Fatalf("Поиск: получено %v, ожидается %v", doc.ID, expID)
	}

	// Проверяем поиск несуществующего документа
	expID = 88
	doc, found = tree.FindDoc(expID)
	if found {
		t.Fatalf("Поиск %v: найден несуществующий документ", expID)
	}

	json, err := tree.JSONData()
	if err != nil {
		t.Fatalf("Ошибка сериализации в json: %v", err)
	}

	tree2 := New()

	err = tree2.LoadFromJSON(json)
	if err != nil {
		t.Fatalf("Получена ошибка загрузки из json: %v", err)
	}

	// Проверяем кол-во элементов в новом объекте
	gotCount = tree.Count()
	if gotCount != expCount {
		t.Fatalf("Count(3): получено %d, ожидается %d", gotCount, expCount)
	}

	// Проверяем поиск по новому объекту
	expID = 8
	doc, found = tree2.FindDoc(expID)
	if !found {
		t.Fatalf("Поиск %v: документ не найден", expID)
	}
	if doc.ID != expID {
		t.Fatalf("Поиск: получено %v, ожидается %v", doc.ID, expID)
	}

	//tree3 := New()

	err = tree2.LoadFromJSON([]byte("не json строка"))
	if err == nil {
		t.Fatalf("Ожидается ошибка загрузки из json, получили nil")
	}

	// Проверяем кол-во элементов в новом объекте, структура не должна измениться
	expCount = 6
	gotCount = tree2.Count()
	if gotCount != expCount {
		t.Fatalf("Count(4): получено %d, ожидается %d", gotCount, expCount)
	}

	// Генерация ID
	for i := 0; i < 100; i++ {
		ID, err := tree.GenerateID()
		if err != nil {
			t.Fatal(err)
		}
		if ID == 0 {
			t.Fatal("Ошибка генерации ID")
		}
	}
}

func TestTree_DeleteDoc(t *testing.T) {

	// Тестируемая структура:
	//            10
	//          /   \
	//         8     12
	//        /     /  \
	//       5    11    18

	// Структура тестов
	// id - ID элементов, которые нужно удалить из начального состояния перед проверкой
	tests := []struct {
		name string
		docs []crawler.Document
		ids  []int
		want []int
	}{
		{
			name: "Начальное состояние",
			docs: docs,
			ids:  []int{},
			want: []int{5, 8, 10, 11, 12, 18},
		},
		{
			name: "Удаление узла без детей (5)",
			docs: docs,
			ids:  []int{5},
			want: []int{8, 10, 11, 12, 18},
		},
		{
			name: "Удаление узла только с левым дочерним (8)",
			docs: docs,
			ids:  []int{8},
			want: []int{5, 10, 11, 12, 18},
		},
		{
			name: "Удаление узла только с правым дочерним (11, 12)",
			docs: docs,
			ids:  []int{11, 12},
			want: []int{5, 8, 10, 18},
		},
		{
			name: "Удаление узла с обоими детьми (12)",
			docs: docs,
			ids:  []int{12},
			want: []int{5, 8, 10, 11, 18},
		},
		{
			name: "Удаление корневого узла (10)",
			docs: docs,
			ids:  []int{10},
			want: []int{5, 8, 11, 12, 18},
		},
		{
			name: "Удаление узла минимального узла в правом поддереве (11)",
			docs: docs,
			ids:  []int{11},
			want: []int{5, 8, 10, 12, 18},
		},
		{
			name: "Удаление самого левого узла (18)",
			docs: docs,
			ids:  []int{18},
			want: []int{5, 8, 10, 11, 12},
		},
		// Тесты специфических кейсов (на которых ловил ошибки)
		{
			name: "Удаление из пустого дерева",
			docs: []crawler.Document{},
			ids:  []int{},
			want: []int{},
		},
		{
			name: "Удаление из дерева с одним узлом",
			docs: []crawler.Document{{ID: 1}},
			ids:  []int{1},
			want: []int{},
		},
		{
			name: "Удаление корня из дерева с двумя узлами",
			docs: []crawler.Document{{ID: 1}, {ID: 2}},
			ids:  []int{1},
			want: []int{2},
		},
		{
			name: "Удаление корня из дерева с левым и правым узлами",
			docs: []crawler.Document{{ID: 2}, {ID: 1}, {ID: 3}},
			ids:  []int{1},
			want: []int{2, 3},
		},
		{
			name: "Удаление левого узла из дерева с двумя узлами",
			docs: []crawler.Document{{ID: 2}, {ID: 1}},
			ids:  []int{1},
			want: []int{2},
		},
		{
			name: "Удаление правого узла из дерева с двумя узлами",
			docs: []crawler.Document{{ID: 1}, {ID: 2}},
			ids:  []int{1},
			want: []int{2},
		},
	}
	for _, tt := range tests {
		tree := New()
		for _, doc := range tt.docs {
			tree.AddDoc(doc)
		}
		for _, id := range tt.ids {
			if id > 0 {
				tree.DeleteDoc(id)
			}
		}
		got := allPlainIDs(tree)
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("%v: получено %v, ожидается %v", tt.name, got, tt.want)
		}
	}
}

// Helper для получения списка IDs всех документов
func allPlainIDs(tree *Tree) []int {
	result := []int{}
	for _, doc := range tree.All() {
		result = append(result, doc.ID)
	}
	return result
}
