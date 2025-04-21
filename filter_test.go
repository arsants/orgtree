package orgtree

import (
	"regexp"
	"testing"
)

func TestFilterSubtreeByRegex(t *testing.T) {
	// Создаем тестовое дерево
	root := NewNode("root")
	child1 := NewNode("test1")
	child2 := NewNode("test2")
	grandchild1 := NewNode("test3")
	grandchild2 := NewNode("other")

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild1)
	child2.AddChild(grandchild2)

	// Тестируем фильтрацию по регулярному выражению
	pattern := regexp.MustCompile("test")
	filteredTree, err := root.FilterSubtreeByRegex(pattern.String())
	if err != nil {
		t.Fatalf("Failed to filter subtree: %v", err)
	}

	// Проверяем структуру отфильтрованного дерева
	if filteredTree.Value != "root" {
		t.Errorf("Expected root value 'root', got '%v'", filteredTree.Value)
	}

	if len(filteredTree.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(filteredTree.Children))
	}

	// Проверяем первого ребенка
	child1Filtered := filteredTree.Children[0]
	if child1Filtered.Value != "test1" {
		t.Errorf("Expected child1 value 'test1', got '%v'", child1Filtered.Value)
	}

	if len(child1Filtered.Children) != 1 {
		t.Errorf("Expected 1 grandchild for child1, got %d", len(child1Filtered.Children))
	}

	// Проверяем второго ребенка
	child2Filtered := filteredTree.Children[1]
	if child2Filtered.Value != "test2" {
		t.Errorf("Expected child2 value 'test2', got '%v'", child2Filtered.Value)
	}

	if len(child2Filtered.Children) != 0 {
		t.Errorf("Expected 0 grandchildren for child2, got %d", len(child2Filtered.Children))
	}

	// Проверяем внука
	grandchild1Filtered := child1Filtered.Children[0]
	if grandchild1Filtered.Value != "test3" {
		t.Errorf("Expected grandchild1 value 'test3', got '%v'", grandchild1Filtered.Value)
	}

	if len(grandchild1Filtered.Children) != 0 {
		t.Errorf("Expected 0 children for grandchild1, got %d", len(grandchild1Filtered.Children))
	}
}

func TestFilterSubtree(t *testing.T) {
	// Создаем тестовое дерево
	root := NewNode("root")
	child1 := NewNode("test1")
	child2 := NewNode("test2")
	grandchild1 := NewNode("test3")
	grandchild2 := NewNode("other")

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild1)
	child2.AddChild(grandchild2)

	// Тестируем фильтрацию по условию
	predicate := func(value interface{}) bool {
		if str, ok := value.(string); ok {
			return str == "test1" || str == "test3"
		}
		return false
	}

	filteredTree := root.FilterSubtree(predicate)

	// Проверяем структуру отфильтрованного дерева
	if filteredTree.Value != "root" {
		t.Errorf("Expected root value 'root', got '%v'", filteredTree.Value)
	}

	if len(filteredTree.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(filteredTree.Children))
	}

	// Проверяем первого ребенка
	child1Filtered := filteredTree.Children[0]
	if child1Filtered.Value != "test1" {
		t.Errorf("Expected child1 value 'test1', got '%v'", child1Filtered.Value)
	}

	if len(child1Filtered.Children) != 1 {
		t.Errorf("Expected 1 grandchild for child1, got %d", len(child1Filtered.Children))
	}

	// Проверяем внука
	grandchild1Filtered := child1Filtered.Children[0]
	if grandchild1Filtered.Value != "test3" {
		t.Errorf("Expected grandchild1 value 'test3', got '%v'", grandchild1Filtered.Value)
	}

	if len(grandchild1Filtered.Children) != 0 {
		t.Errorf("Expected 0 children for grandchild1, got %d", len(grandchild1Filtered.Children))
	}
}
