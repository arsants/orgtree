package orgtree

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"hash"
)

// Node представляет собой узел дерева с динамическими данными
type Node struct {
	Value    interface{} `json:"value"`
	Children []*Node     `json:"children"`
}

// NewNode создает новый узел
func NewNode(value interface{}) *Node {
	return &Node{
		Value:    value,
		Children: []*Node{},
	}
}

// AddChild добавляет дочерний узел
func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

// Find ищет первый узел с данным значением
func (n *Node) Find(value interface{}) *Node {
	if n.Value == value {
		return n
	}
	for _, child := range n.Children {
		if found := child.Find(value); found != nil {
			return found
		}
	}
	return nil
}

// SubTree возвращает поддерево с корнем в найденном узле
func (n *Node) SubTree(value interface{}) *Node {
	return n.Find(value)
}

// ToJSON сериализует дерево в JSON
func (n *Node) ToJSON() ([]byte, error) {
	return json.MarshalIndent(n, "", "  ")
}

// FromJSON десериализует JSON в дерево
func FromJSON(data []byte) (*Node, error) {
	var root Node
	err := json.Unmarshal(data, &root)
	return &root, err
}

// Hash возвращает хеш дерева, учитывая значения всех узлов
func (n *Node) Hash() []byte {
	h := sha256.New()
	n.hashHelper(h)
	return h.Sum(nil)
}

func (n *Node) hashHelper(h hash.Hash) {
	valBytes, _ := json.Marshal(n.Value)
	h.Write(valBytes)
	for _, child := range n.Children {
		child.hashHelper(h)
	}
}

// HashString возвращает хеш дерева в виде hex строки
func (n *Node) HashString() string {
	return fmt.Sprintf("%x", n.Hash())
}
