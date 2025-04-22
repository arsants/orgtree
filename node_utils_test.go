package orgtree

import (
	"bytes"
	"os"
	"testing"
)

func TestPrintTree(t *testing.T) {
	// Создаем тестовое дерево
	root := NewNode("root")
	child1 := NewNode("child1")
	child2 := NewNode("child2")
	grandchild1 := NewNode("grandchild1")
	grandchild2 := NewNode("grandchild2")

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild1)
	child2.AddChild(grandchild2)

	// Перенаправляем вывод в буфер
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Выводим дерево
	root.PrintTree()

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = old

	// Читаем вывод
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Проверяем структуру вывода
	expected := `└── root
    ├── child1
    │   └── grandchild1
    └── child2
        └── grandchild2
`

	if output != expected {
		t.Errorf("PrintTree output mismatch.\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestWalkTree(t *testing.T) {
	// Создаем тестовое дерево
	root := NewNode("root")
	child1 := NewNode("child1")
	child2 := NewNode("child2")
	grandchild1 := NewNode("grandchild1")
	grandchild2 := NewNode("grandchild2")

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild1)
	child2.AddChild(grandchild2)

	// Собираем информацию о посещенных узлах
	visited := make(map[string]int)
	root.WalkTree(func(node *Node, depth int) {
		if str, ok := node.Value.(string); ok {
			visited[str] = depth
		}
	})

	// Проверяем глубину каждого узла
	expected := map[string]int{
		"root":        0,
		"child1":      1,
		"child2":      1,
		"grandchild1": 2,
		"grandchild2": 2,
	}

	for node, depth := range expected {
		if visited[node] != depth {
			t.Errorf("Node %s: expected depth %d, got %d", node, depth, visited[node])
		}
	}
}

func TestFind(t *testing.T) {
	// Создаем тестовое дерево
	root := NewNode("root")
	child1 := NewNode("child1")
	child2 := NewNode("child2")
	grandchild1 := NewNode("grandchild1")
	grandchild2 := NewNode("grandchild2")

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild1)
	child2.AddChild(grandchild2)

	// Тест 1: Поиск существующего узла
	if node := root.Find("child1"); node == nil {
		t.Error("Find failed to find existing node 'child1'")
	} else if node.Value != "child1" {
		t.Errorf("Find returned wrong node. Expected 'child1', got '%v'", node.Value)
	}

	// Тест 2: Поиск несуществующего узла
	if node := root.Find("nonexistent"); node != nil {
		t.Error("Find returned node for non-existent value")
	}

	// Тест 3: Поиск узла в поддереве
	if node := root.Find("grandchild1"); node == nil {
		t.Error("Find failed to find node 'grandchild1' in subtree")
	} else if node.Value != "grandchild1" {
		t.Errorf("Find returned wrong node. Expected 'grandchild1', got '%v'", node.Value)
	}
}

func TestFilter(t *testing.T) {
	// Создаем тестовое дерево
	root := NewNode("root")
	child1 := NewNode("child1")
	child2 := NewNode("child2")
	grandchild1 := NewNode("grandchild1")
	grandchild2 := NewNode("grandchild2")

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild1)
	child2.AddChild(grandchild2)

	// Тест 1: Фильтрация по точному совпадению
	predicate1 := func(value interface{}) bool {
		if str, ok := value.(string); ok {
			return str == "child1"
		}
		return false
	}

	if node := root.Filter(predicate1); node == nil {
		t.Error("Filter failed to find node matching predicate")
	} else if node.Value != "child1" {
		t.Errorf("Filter returned wrong node. Expected 'child1', got '%v'", node.Value)
	}

	// Тест 2: Фильтрация по условию
	predicate2 := func(value interface{}) bool {
		if str, ok := value.(string); ok {
			return len(str) > 5
		}
		return false
	}

	if node := root.Filter(predicate2); node == nil {
		t.Error("Filter failed to find node matching length predicate")
	} else if str, ok := node.Value.(string); !ok || len(str) <= 5 {
		t.Error("Filter returned node that doesn't match length predicate")
	}
}

func TestSubTree(t *testing.T) {
	// Создаем тестовое дерево
	root := NewNode("root")
	child1 := NewNode("child1")
	child2 := NewNode("child2")
	grandchild1 := NewNode("grandchild1")
	grandchild2 := NewNode("grandchild2")

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild1)
	child2.AddChild(grandchild2)

	// Тест 1: Получение поддерева от существующего узла
	subtree := root.SubTree("child1")
	if subtree == nil {
		t.Error("SubTree failed to find existing subtree")
	} else if subtree.Value != "child1" {
		t.Errorf("SubTree returned wrong root. Expected 'child1', got '%v'", subtree.Value)
	} else if len(subtree.Children) != 1 {
		t.Errorf("SubTree has wrong number of children. Expected 1, got %d", len(subtree.Children))
	}

	// Тест 2: Получение поддерева от несуществующего узла
	if subtree := root.SubTree("nonexistent"); subtree != nil {
		t.Error("SubTree returned non-nil for non-existent value")
	}

	// Тест 3: Проверка структуры поддерева
	subtree = root.SubTree("child2")
	if subtree == nil {
		t.Error("SubTree failed to find subtree for 'child2'")
	} else {
		if len(subtree.Children) != 1 {
			t.Errorf("SubTree has wrong number of children. Expected 1, got %d", len(subtree.Children))
		}
		if subtree.Children[0].Value != "grandchild2" {
			t.Errorf("SubTree has wrong child. Expected 'grandchild2', got '%v'", subtree.Children[0].Value)
		}
	}
}

func TestHashString(t *testing.T) {
	// Создаем тестовое дерево
	root := NewNode("root")
	child1 := NewNode("child1")
	child2 := NewNode("child2")
	grandchild1 := NewNode("grandchild1")
	grandchild2 := NewNode("grandchild2")

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild1)
	child2.AddChild(grandchild2)

	// Получаем хеш
	hash1 := root.HashString()

	// Создаем идентичное дерево
	root2 := NewNode("root")
	child1_2 := NewNode("child1")
	child2_2 := NewNode("child2")
	grandchild1_2 := NewNode("grandchild1")
	grandchild2_2 := NewNode("grandchild2")

	root2.AddChild(child1_2)
	root2.AddChild(child2_2)
	child1_2.AddChild(grandchild1_2)
	child2_2.AddChild(grandchild2_2)

	// Получаем хеш второго дерева
	hash2 := root2.HashString()

	// Проверяем, что хеши идентичных деревьев совпадают
	if hash1 != hash2 {
		t.Errorf("Expected identical trees to have identical hashes, got %s and %s", hash1, hash2)
	}

	// Изменяем значение в одном из узлов
	child1.Value = "changed"

	// Получаем новый хеш
	hash3 := root.HashString()

	// Проверяем, что хеш изменился
	if hash1 == hash3 {
		t.Error("Expected hash to change after modifying node value")
	}

	// Создаем дерево с другой структурой
	root4 := NewNode("root")
	child1_4 := NewNode("child1")
	child2_4 := NewNode("child2")
	grandchild1_4 := NewNode("grandchild1")

	root4.AddChild(child1_4)
	root4.AddChild(child2_4)
	child1_4.AddChild(grandchild1_4)

	// Получаем хеш дерева с другой структурой
	hash4 := root4.HashString()

	// Проверяем, что хеш отличается
	if hash1 == hash4 {
		t.Error("Expected different tree structures to have different hashes")
	}
}
