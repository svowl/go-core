package btree

import (
	"gosearch/pkg/crawler"
	"testing"
)

// Входные данные для тестирования
var testDocs = []crawler.Document{
	{ID: 10, URL: "https://google.com", Title: "Google"},
	{ID: 12, URL: "https://go.dev/", Title: "Why Go"},
	{ID: 8, URL: "https://golang.org/", Title: "The Go Programming Language"},
	{ID: 18, URL: "https://yandex.ru/", Title: "Yandex"},
	{ID: 5, URL: "https://rambler.ru/", Title: "Rambler"},
	{ID: 11, URL: "https://mail.ru/", Title: "Mail.ru"},
	{ID: 5, URL: "https://yahoo.com/", Title: "Yahoo!"},
	{ID: 0, URL: "https://yahoo.com/", Title: "Yahoo!"},
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
	tree.AddDoc(testDocs[0])
	if tree.Root.Doc.ID != testDocs[0].ID {
		t.Fatalf("Добавление корневого узла %v: добавлен %v", testDocs[0].ID, tree.Root.Doc.ID)
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
	tree.AddDoc(testDocs[1])
	if tree.Root.Right.Doc.ID != testDocs[1].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", testDocs[1].ID, tree.Root.Right.Doc.ID)
	}

	// Добавляем узел 8:  10
	//                   /   \
	//                  8     12
	tree.AddDoc(testDocs[2])
	if tree.Root.Left.Doc.ID != testDocs[2].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", testDocs[2].ID, tree.Root.Left.Doc.ID)
	}

	// Добавляем узел 18:  10
	//                    /   \
	//                   8     12
	//                           \
	//                            18
	tree.AddDoc(testDocs[3])
	if tree.Root.Right.Right.Doc.ID != testDocs[3].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", testDocs[3].ID, tree.Root.Right.Right.Doc.ID)
	}

	// Добавляем узел 5:   10
	//                    /   \
	//                   8     12
	//                  /        \
	//                 5          18
	tree.AddDoc(testDocs[4])
	if tree.Root.Left.Left.Doc.ID != testDocs[4].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", testDocs[4].ID, tree.Root.Left.Left.Doc.ID)
	}

	// Добавляем узел 11:   10
	//                    /   \
	//                   8     12
	//                  /     /  \
	//                 5    11    18
	tree.AddDoc(testDocs[5])
	if tree.Root.Right.Left.Doc.ID != testDocs[5].ID {
		t.Fatalf("Добавление узла %v: добавлен %v", testDocs[5].ID, tree.Root.Right.Left.Doc.ID)
	}

	// В дереве 6 элементов, проверяем Count
	expCount = 6
	gotCount = tree.Count()
	if gotCount != expCount {
		t.Fatalf("Count(1): получено %d, ожидается %d", gotCount, expCount)
	}

	// Добавляем документ с существующим ID (5), документ должен переписаться
	tree.AddDoc(testDocs[6])
	if tree.Root.Left.Left.Doc.Title != testDocs[6].Title {
		t.Fatalf("Добавление узла %v: добавлен %v", testDocs[6].Title, tree.Root.Left.Left.Doc.Title)
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
