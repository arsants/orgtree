package orgtree

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	node := NewNode("test")
	if node.Value != "test" {
		t.Errorf("Expected value 'test', got '%v'", node.Value)
	}
	if len(node.Children) != 0 {
		t.Error("Expected empty children slice")
	}
}

func TestAddChild(t *testing.T) {
	parent := NewNode("parent")
	child := NewNode("child")

	parent.AddChild(child)

	if len(parent.Children) != 1 {
		t.Error("Expected one child")
	}
	if parent.Children[0] != child {
		t.Error("Child not properly added")
	}
}

func TestGetPath(t *testing.T) {
	root := NewNode("root")
	child := NewNode("child")
	grandchild := NewNode("grandchild")

	root.AddChild(child)
	child.AddChild(grandchild)

	path := grandchild.GetPath(root)
	expected := []string{"root", "child", "grandchild"}

	if len(path) != len(expected) {
		t.Errorf("Expected path length %d, got %d", len(expected), len(path))
	}

	for i, node := range path {
		if str, ok := node.Value.(string); !ok || str != expected[i] {
			t.Errorf("Expected path[%d] = %s, got %v", i, expected[i], node.Value)
		}
	}
}

func TestGetDepth(t *testing.T) {
	root := NewNode("root")
	child := NewNode("child")
	grandchild := NewNode("grandchild")

	root.AddChild(child)
	child.AddChild(grandchild)

	tests := []struct {
		node     *Node
		expected int
	}{
		{root, 0},
		{child, 1},
		{grandchild, 2},
	}

	for _, test := range tests {
		depth, ok := test.node.GetDepth(root)
		if !ok {
			t.Errorf("Failed to get depth for node %v", test.node.Value)
		}
		if depth != test.expected {
			t.Errorf("For node %v, expected depth %d, got %d",
				test.node.Value, test.expected, depth)
		}
	}
}

func TestTreeJSONSerialization(t *testing.T) {
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

	// Сериализуем дерево
	jsonData, err := root.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize tree: %v", err)
	}

	// Десериализуем обратно
	newRoot, err := FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to deserialize tree: %v", err)
	}

	// Проверяем, что все узлы на месте с помощью итератора
	originalNodes := make(map[string]bool)
	it := NewPreOrderIterator(root)
	for node := it.Next(); node != nil; node = it.Next() {
		if str, ok := node.Value.(string); ok {
			originalNodes[str] = true
		}
	}

	newNodes := make(map[string]bool)
	it = NewPreOrderIterator(newRoot)
	for node := it.Next(); node != nil; node = it.Next() {
		if str, ok := node.Value.(string); ok {
			newNodes[str] = true
		}
	}

	// Проверяем, что все узлы сохранились
	for value := range originalNodes {
		if !newNodes[value] {
			t.Errorf("Node '%s' lost after serialization/deserialization", value)
		}
	}

	for value := range newNodes {
		if !originalNodes[value] {
			t.Errorf("Extra node '%s' appeared after serialization/deserialization", value)
		}
	}
}
