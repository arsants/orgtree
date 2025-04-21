package orgtree

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"hash"
)

// GetDepth возвращает глубину текущего узла
func (n *Node) GetDepth(root *Node) (int, bool) {
	it := NewPreOrderIterator(root)
	for current := it.Next(); current != nil; current = it.Next() {
		if current == n {
			return len(current.GetPath(root)) - 1, true
		}
	}
	return 0, false
}

// GetPath возвращает путь от корня до текущего узла
func (n *Node) GetPath(root *Node) []*Node {
	path := []*Node{}
	var dfs func(node *Node) bool
	dfs = func(node *Node) bool {
		path = append(path, node)
		if node == n {
			return true
		}
		for _, child := range node.Children {
			if dfs(child) {
				return true
			}
		}
		path = path[:len(path)-1]
		return false
	}
	if dfs(root) {
		return path
	}
	return []*Node{}
}

// Find ищет первый узел с данным значением
func (n *Node) Find(value interface{}) *Node {
	it := NewPreOrderIterator(n)
	for node := it.Next(); node != nil; node = it.Next() {
		if node.Value == value {
			return node
		}
	}
	return nil
}

// Filter возвращает первый узел, удовлетворяющий предикату
func (n *Node) Filter(predicate func(interface{}) bool) *Node {
	it := NewPreOrderIterator(n)
	for node := it.Next(); node != nil; node = it.Next() {
		if predicate(node.Value) {
			return node
		}
	}
	return nil
}

// SubTree возвращает поддерево с корнем в найденном узле
func (n *Node) SubTree(value interface{}) *Node {
	return n.Find(value) // уже дерево от нужного корня
}

// ToJSON сериализует дерево в JSON
func (n *Node) ToJSON() ([]byte, error) {
	return json.MarshalIndent(n, "", "  ")
}

// FromJSON десериализует JSON в дерево
func FromJSON(data []byte) (*Node, error) {
	var root Node
	err := json.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}
	return &root, nil
}

// Hash возвращает хеш дерева, учитывая значения всех узлов
func (n *Node) Hash() []byte {
	h := sha256.New()
	n.hashHelper(h)
	return h.Sum(nil)
}

// HashString возвращает строковое представление хеша
func (n *Node) HashString() string {
	return fmt.Sprintf("%x", n.Hash())
}

func (n *Node) hashHelper(h hash.Hash) {
	valBytes, _ := json.Marshal(n.Value)
	h.Write(valBytes)
	for _, child := range n.Children {
		child.hashHelper(h)
	}
}
