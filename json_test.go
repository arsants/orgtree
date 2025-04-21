package orgtree

import (
	"encoding/json"
	"testing"
)

func TestToJSON(t *testing.T) {
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

	// Сериализуем в JSON
	jsonData, err := root.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize to JSON: %v", err)
	}

	// Проверяем, что JSON валидный
	var jsonTree map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonTree); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Проверяем структуру JSON
	if jsonTree["value"] != "root" {
		t.Errorf("Expected root value 'root', got '%v'", jsonTree["value"])
	}

	children, ok := jsonTree["children"].([]interface{})
	if !ok {
		t.Fatal("Expected children to be an array")
	}

	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}

	// Проверяем первого ребенка
	child1Map, ok := children[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected child1 to be an object")
	}

	if child1Map["value"] != "child1" {
		t.Errorf("Expected child1 value 'child1', got '%v'", child1Map["value"])
	}

	// Проверяем второго ребенка
	child2Map, ok := children[1].(map[string]interface{})
	if !ok {
		t.Fatal("Expected child2 to be an object")
	}

	if child2Map["value"] != "child2" {
		t.Errorf("Expected child2 value 'child2', got '%v'", child2Map["value"])
	}
}

func TestFromJSON(t *testing.T) {
	// Создаем тестовый JSON
	jsonStr := `{
		"value": "root",
		"children": [
			{
				"value": "child1",
				"children": [
					{
						"value": "grandchild1",
						"children": []
					}
				]
			},
			{
				"value": "child2",
				"children": [
					{
						"value": "grandchild2",
						"children": []
					}
				]
			}
		]
	}`

	// Десериализуем из JSON
	root, err := FromJSON([]byte(jsonStr))
	if err != nil {
		t.Fatalf("Failed to deserialize from JSON: %v", err)
	}

	// Проверяем структуру дерева
	if root.Value != "root" {
		t.Errorf("Expected root value 'root', got '%v'", root.Value)
	}

	if len(root.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(root.Children))
	}

	// Проверяем первого ребенка
	child1 := root.Children[0]
	if child1.Value != "child1" {
		t.Errorf("Expected child1 value 'child1', got '%v'", child1.Value)
	}

	if len(child1.Children) != 1 {
		t.Errorf("Expected 1 grandchild for child1, got %d", len(child1.Children))
	}

	// Проверяем второго ребенка
	child2 := root.Children[1]
	if child2.Value != "child2" {
		t.Errorf("Expected child2 value 'child2', got '%v'", child2.Value)
	}

	if len(child2.Children) != 1 {
		t.Errorf("Expected 1 grandchild for child2, got %d", len(child2.Children))
	}
}
