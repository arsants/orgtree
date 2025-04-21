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

func TestFind(t *testing.T) {
	root := NewNode("root")
	child1 := NewNode("child1")
	child2 := NewNode("child2")

	root.AddChild(child1)
	root.AddChild(child2)

	found := root.Find("child1")
	if found == nil {
		t.Error("Expected to find node with value 'child1'")
	}
	if found.Value != "child1" {
		t.Errorf("Expected value 'child1', got '%v'", found.Value)
	}
}

func TestToJSON(t *testing.T) {
	root := NewNode("root")
	child := NewNode("child")
	root.AddChild(child)

	data, err := root.ToJSON()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Error("Expected non-empty JSON data")
	}
}

func TestFromJSON(t *testing.T) {
	root := NewNode("root")
	child := NewNode("child")
	root.AddChild(child)

	data, _ := root.ToJSON()
	loaded, err := FromJSON(data)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if loaded.Value != "root" {
		t.Errorf("Expected value 'root', got '%v'", loaded.Value)
	}
	if len(loaded.Children) != 1 {
		t.Error("Expected one child after loading from JSON")
	}
}

func TestHash(t *testing.T) {
	root := NewNode("root")
	child := NewNode("child")
	root.AddChild(child)

	hash := root.Hash()
	if len(hash) == 0 {
		t.Error("Expected non-empty hash")
	}

	hashStr := root.HashString()
	if len(hashStr) == 0 {
		t.Error("Expected non-empty hash string")
	}
}
