package orgtree

import (
	"testing"
)

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
