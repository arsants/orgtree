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

// ToJSON сериализует дерево в JSON, гарантируя включение всех дочерних элементов
func (n *Node) ToJSON() ([]byte, error) {
	// Создаем карту для отслеживания обработанных узлов
	processed := make(map[*Node]bool)

	// Функция для конвертации Node в map
	var convertNode func(*Node) map[string]interface{}
	convertNode = func(node *Node) map[string]interface{} {
		if processed[node] {
			return nil
		}
		processed[node] = true

		children := make([]interface{}, 0, len(node.Children))
		for _, child := range node.Children {
			if childMap := convertNode(child); childMap != nil {
				children = append(children, childMap)
			}
		}

		return map[string]interface{}{
			"value":    node.Value,
			"children": children,
		}
	}

	// Используем BFS для проверки, что все узлы будут обработаны
	it := NewBFSIterator(n)
	for node := it.Next(); node != nil; node = it.Next() {
		processed[node] = false
	}

	// Конвертируем дерево
	treeMap := convertNode(n)

	// Сериализуем в JSON
	return json.MarshalIndent(treeMap, "", "  ")
}

// FromJSON десериализует JSON в дерево
func FromJSON(data []byte) (*Node, error) {
	var treeMap map[string]interface{}
	if err := json.Unmarshal(data, &treeMap); err != nil {
		return nil, err
	}

	// Функция для конвертации map в Node
	var convertMap func(map[string]interface{}) *Node
	convertMap = func(m map[string]interface{}) *Node {
		node := NewNode(m["value"])
		if children, ok := m["children"].([]interface{}); ok {
			for _, child := range children {
				if childMap, ok := child.(map[string]interface{}); ok {
					node.AddChild(convertMap(childMap))
				}
			}
		}
		return node
	}

	return convertMap(treeMap), nil
}

// PrintTree выводит дерево в консоль с отступами
func (n *Node) PrintTree() {
	n.printTreeRecursive("", true)
}

// printTreeRecursive рекурсивно выводит дерево с отступами
func (n *Node) printTreeRecursive(prefix string, isLast bool) {
	// Определяем символы для отображения структуры
	var marker string
	if isLast {
		marker = "└── "
	} else {
		marker = "├── "
	}

	// Выводим текущий узел
	fmt.Printf("%s%s%v\n", prefix, marker, n.Value)

	// Определяем префикс для детей
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Рекурсивно обходим детей
	for i, child := range n.Children {
		isLastChild := i == len(n.Children)-1
		child.printTreeRecursive(childPrefix, isLastChild)
	}
}

// WalkTree обходит дерево рекурсивно и вызывает callback для каждого узла
func (n *Node) WalkTree(callback func(*Node, int)) {
	n.walkTreeRecursive(callback, 0)
}

// walkTreeRecursive выполняет рекурсивный обход дерева
func (n *Node) walkTreeRecursive(callback func(*Node, int), depth int) {
	// Вызываем callback для текущего узла
	callback(n, depth)

	// Рекурсивно обходим всех детей
	for _, child := range n.Children {
		child.walkTreeRecursive(callback, depth+1)
	}
}
